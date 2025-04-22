package consumable

import (
	consumer2 "atlas-channel/kafka/consumer"
	consumable2 "atlas-channel/kafka/message/consumable"
	_map "atlas-channel/map"
	"atlas-channel/server"
	"atlas-channel/session"
	model2 "atlas-channel/socket/model"
	"atlas-channel/socket/writer"
	"context"
	"github.com/Chronicle20/atlas-kafka/consumer"
	"github.com/Chronicle20/atlas-kafka/handler"
	"github.com/Chronicle20/atlas-kafka/message"
	"github.com/Chronicle20/atlas-kafka/topic"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

func InitConsumers(l logrus.FieldLogger) func(func(config consumer.Config, decorators ...model.Decorator[consumer.Config])) func(consumerGroupId string) {
	return func(rf func(config consumer.Config, decorators ...model.Decorator[consumer.Config])) func(consumerGroupId string) {
		return func(consumerGroupId string) {
			rf(consumer2.NewConfig(l)("consumable_command")(consumable2.EnvEventTopic)(consumerGroupId), consumer.SetHeaderParsers(consumer.SpanHeaderParser, consumer.TenantHeaderParser), consumer.SetStartOffset(kafka.LastOffset))
		}
	}
}

func InitHandlers(l logrus.FieldLogger) func(sc server.Model) func(wp writer.Producer) func(rf func(topic string, handler handler.Handler) (string, error)) {
	return func(sc server.Model) func(wp writer.Producer) func(rf func(topic string, handler handler.Handler) (string, error)) {
		return func(wp writer.Producer) func(rf func(topic string, handler handler.Handler) (string, error)) {
			return func(rf func(topic string, handler handler.Handler) (string, error)) {
				var t string
				t, _ = topic.EnvProvider(l)(consumable2.EnvEventTopic)()
				_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleErrorConsumableEvent(sc, wp))))
				_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleScrollConsumableEvent(sc, wp))))
			}
		}
	}
}

func handleErrorConsumableEvent(sc server.Model, wp writer.Producer) message.Handler[consumable2.Event[consumable2.ErrorBody]] {
	return func(l logrus.FieldLogger, ctx context.Context, e consumable2.Event[consumable2.ErrorBody]) {
		if e.Type != consumable2.EventTypeError {
			return
		}

		if e.Body.Error == consumable2.ErrorTypePetCannotConsume {
			err := session.NewProcessor(l, ctx).IfPresentByCharacterId(sc.WorldId(), sc.ChannelId())(e.CharacterId, session.Announce(l)(ctx)(wp)(writer.PetCashFoodResult)(writer.PetCashFoodErrorResultBody()))
			if err != nil {
				l.WithError(err).Errorf("Unable to process error event for character [%d].", e.CharacterId)
			}
			return
		}

		err := session.NewProcessor(l, ctx).IfPresentByCharacterId(sc.WorldId(), sc.ChannelId())(e.CharacterId, session.Announce(l)(ctx)(wp)(writer.StatChanged)(writer.StatChangedBody(l)(make([]model2.StatUpdate, 0), true)))
		if err != nil {
			l.WithError(err).Errorf("Unable to process error event for character [%d].", e.CharacterId)
		}
	}
}

func handleScrollConsumableEvent(sc server.Model, wp writer.Producer) message.Handler[consumable2.Event[consumable2.ScrollBody]] {
	return func(l logrus.FieldLogger, ctx context.Context, e consumable2.Event[consumable2.ScrollBody]) {
		if e.Type != consumable2.EventTypeScroll {
			return
		}

		err := session.NewProcessor(l, ctx).IfPresentByCharacterId(sc.WorldId(), sc.ChannelId())(e.CharacterId, func(s session.Model) error {
			return _map.NewProcessor(l, ctx).ForSessionsInMap(s.Map(), session.Announce(l)(ctx)(wp)(writer.CharacterItemUpgrade)(writer.CharacterItemUpgradeBody(e.CharacterId, e.Body.Success, e.Body.Cursed, e.Body.LegendarySpirit, e.Body.WhiteScroll)))
		})
		if err != nil {
			l.WithError(err).Errorf("Unable to process scroll event for character [%d].", e.CharacterId)
		}
	}
}
