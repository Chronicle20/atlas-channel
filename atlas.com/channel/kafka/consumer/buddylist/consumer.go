package buddylist

import (
	"atlas-channel/buddylist"
	consumer2 "atlas-channel/kafka/consumer"
	"atlas-channel/server"
	"atlas-channel/session"
	"atlas-channel/socket/writer"
	"context"
	"github.com/Chronicle20/atlas-constants/world"
	"github.com/Chronicle20/atlas-kafka/consumer"
	"github.com/Chronicle20/atlas-kafka/handler"
	"github.com/Chronicle20/atlas-kafka/message"
	"github.com/Chronicle20/atlas-kafka/topic"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/Chronicle20/atlas-tenant"
	"github.com/sirupsen/logrus"
)

func InitConsumers(l logrus.FieldLogger) func(func(config consumer.Config, decorators ...model.Decorator[consumer.Config])) func(consumerGroupId string) {
	return func(rf func(config consumer.Config, decorators ...model.Decorator[consumer.Config])) func(consumerGroupId string) {
		return func(consumerGroupId string) {
			rf(consumer2.NewConfig(l)("buddy_list_status_event")(EnvStatusEventTopic)(consumerGroupId), consumer.SetHeaderParsers(consumer.SpanHeaderParser, consumer.TenantHeaderParser))
		}
	}
}

func InitHandlers(l logrus.FieldLogger) func(sc server.Model) func(wp writer.Producer) func(rf func(topic string, handler handler.Handler) (string, error)) {
	return func(sc server.Model) func(wp writer.Producer) func(rf func(topic string, handler handler.Handler) (string, error)) {
		return func(wp writer.Producer) func(rf func(topic string, handler handler.Handler) (string, error)) {
			return func(rf func(topic string, handler handler.Handler) (string, error)) {
				var t string
				t, _ = topic.EnvProvider(l)(EnvStatusEventTopic)()
				_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleStatusEventBuddyAdded(sc, wp))))
				_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleStatusEventBuddyRemoved(sc, wp))))
				_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleStatusEventBuddyUpdated(sc, wp))))
				_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleStatusEventBuddyChannelChange(sc, wp))))
				_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleStatusEventBuddyCapacityChange(sc, wp))))
				_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleStatusEventBuddyError(sc, wp))))
			}
		}
	}
}

func handleStatusEventBuddyAdded(sc server.Model, wp writer.Producer) message.Handler[statusEvent[buddyAddedStatusEventBody]] {
	return func(l logrus.FieldLogger, ctx context.Context, c statusEvent[buddyAddedStatusEventBody]) {
		if c.Type != StatusEventTypeBuddyAdded {
			return
		}

		if !sc.IsWorld(tenant.MustFromContext(ctx), world.Id(c.WorldId)) {
			return
		}

		err := session.IfPresentByCharacterId(sc.Tenant(), sc.WorldId(), sc.ChannelId())(c.CharacterId, redrawBuddyList(l)(ctx)(wp)())
		if err != nil {
			l.WithError(err).Errorf("Unable to write character [%d] buddy list.", c.CharacterId)
		}
	}
}

func handleStatusEventBuddyRemoved(sc server.Model, wp writer.Producer) message.Handler[statusEvent[buddyRemovedStatusEventBody]] {
	return func(l logrus.FieldLogger, ctx context.Context, c statusEvent[buddyRemovedStatusEventBody]) {
		if c.Type != StatusEventTypeBuddyRemoved {
			return
		}

		if !sc.IsWorld(tenant.MustFromContext(ctx), world.Id(c.WorldId)) {
			return
		}

		err := session.IfPresentByCharacterId(sc.Tenant(), sc.WorldId(), sc.ChannelId())(c.CharacterId, redrawBuddyList(l)(ctx)(wp)())
		if err != nil {
			l.WithError(err).Errorf("Unable to write character [%d] buddy list.", c.CharacterId)
		}
	}
}

func redrawBuddyList(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func() model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func() model.Operator[session.Model] {
		t := tenant.MustFromContext(ctx)
		return func(wp writer.Producer) func() model.Operator[session.Model] {
			return func() model.Operator[session.Model] {
				return func(s session.Model) error {
					bl, err := buddylist.GetById(l)(ctx)(s.CharacterId())
					if err != nil {
						return err
					}

					err = session.Announce(l)(ctx)(wp)(writer.BuddyOperation)(writer.BuddyListUpdateBody(l, t)(bl.Buddies()))(s)
					if err != nil {
						return err
					}
					return nil
				}
			}
		}
	}
}

func handleStatusEventBuddyUpdated(sc server.Model, wp writer.Producer) message.Handler[statusEvent[buddyUpdatedStatusEventBody]] {
	return func(l logrus.FieldLogger, ctx context.Context, c statusEvent[buddyUpdatedStatusEventBody]) {
		if c.Type != StatusEventTypeBuddyUpdated {
			return
		}

		if !sc.IsWorld(tenant.MustFromContext(ctx), world.Id(c.WorldId)) {
			return
		}

		err := session.IfPresentByCharacterId(sc.Tenant(), sc.WorldId(), sc.ChannelId())(c.CharacterId, updateBuddy(l)(ctx)(wp)(c.Body.CharacterId, c.Body.Group, c.Body.CharacterName, c.Body.ChannelId, c.Body.InShop))
		if err != nil {
			l.WithError(err).Errorf("Unable to announce character [%d] buddy [%d] channel change to [%d].", c.CharacterId, c.Body.CharacterId, c.Body.ChannelId)
		}
	}
}

func updateBuddy(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(characterId uint32, group string, characterName string, channelId int8, inShop bool) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(characterId uint32, group string, characterName string, channelId int8, inShop bool) model.Operator[session.Model] {
		t := tenant.MustFromContext(ctx)
		return func(wp writer.Producer) func(characterId uint32, group string, characterName string, channelId int8, inShop bool) model.Operator[session.Model] {
			return func(characterId uint32, group string, characterName string, channelId int8, inShop bool) model.Operator[session.Model] {
				return session.Announce(l)(ctx)(wp)(writer.BuddyOperation)(writer.BuddyUpdateBody(l, t)(characterId, group, characterName, channelId, inShop))
			}
		}
	}
}

func handleStatusEventBuddyChannelChange(sc server.Model, wp writer.Producer) message.Handler[statusEvent[buddyChannelChangeStatusEventBody]] {
	return func(l logrus.FieldLogger, ctx context.Context, c statusEvent[buddyChannelChangeStatusEventBody]) {
		if c.Type != StatusEventTypeBuddyChannelChange {
			return
		}

		if !sc.IsWorld(tenant.MustFromContext(ctx), world.Id(c.WorldId)) {
			return
		}

		err := session.IfPresentByCharacterId(sc.Tenant(), sc.WorldId(), sc.ChannelId())(c.CharacterId, buddyChannelChange(l)(ctx)(wp)(c.Body.CharacterId, c.Body.ChannelId))
		if err != nil {
			l.WithError(err).Errorf("Unable to announce character [%d] buddy [%d] channel change to [%d].", c.CharacterId, c.Body.CharacterId, c.Body.ChannelId)
		}
	}
}

func buddyChannelChange(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(characterId uint32, channelId int8) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(characterId uint32, channelId int8) model.Operator[session.Model] {
		return func(wp writer.Producer) func(characterId uint32, channelId int8) model.Operator[session.Model] {
			return func(characterId uint32, channelId int8) model.Operator[session.Model] {
				return session.Announce(l)(ctx)(wp)(writer.BuddyOperation)(writer.BuddyChannelChangeBody(l)(characterId, channelId))
			}
		}
	}
}

func handleStatusEventBuddyCapacityChange(sc server.Model, wp writer.Producer) message.Handler[statusEvent[buddyCapacityChangeStatusEventBody]] {
	return func(l logrus.FieldLogger, ctx context.Context, c statusEvent[buddyCapacityChangeStatusEventBody]) {
		if c.Type != StatusEventTypeBuddyCapacityUpdate {
			return
		}

		if !sc.IsWorld(tenant.MustFromContext(ctx), world.Id(c.WorldId)) {
			return
		}

		err := session.IfPresentByCharacterId(sc.Tenant(), sc.WorldId(), sc.ChannelId())(c.CharacterId, buddyCapacityChange(l)(ctx)(wp)(c.Body.Capacity))
		if err != nil {
			l.WithError(err).Errorf("Unable to announce character [%d] buddy list capacity [%d] update.", c.CharacterId, c.Body.Capacity)
		}
	}
}

func buddyCapacityChange(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(capacity byte) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(capacity byte) model.Operator[session.Model] {
		return func(wp writer.Producer) func(capacity byte) model.Operator[session.Model] {
			return func(capacity byte) model.Operator[session.Model] {
				return session.Announce(l)(ctx)(wp)(writer.BuddyOperation)(writer.BuddyCapacityUpdateBody(l)(capacity))
			}
		}
	}
}

func handleStatusEventBuddyError(sc server.Model, wp writer.Producer) message.Handler[statusEvent[errorStatusEventBody]] {
	return func(l logrus.FieldLogger, ctx context.Context, c statusEvent[errorStatusEventBody]) {
		if c.Type != StatusEventTypeError {
			return
		}

		if !sc.IsWorld(tenant.MustFromContext(ctx), world.Id(c.WorldId)) {
			return
		}

		err := session.IfPresentByCharacterId(sc.Tenant(), sc.WorldId(), sc.ChannelId())(c.CharacterId, buddyError(l)(ctx)(wp)(c.Body.Error))
		if err != nil {
			l.WithError(err).Errorf("Unable to announce character [%d] error [%s].", c.CharacterId, c.Body.Error)
		}
	}
}

func buddyError(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(errorCode string) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(errorCode string) model.Operator[session.Model] {
		return func(wp writer.Producer) func(errorCode string) model.Operator[session.Model] {
			return func(errorCode string) model.Operator[session.Model] {
				return session.Announce(l)(ctx)(wp)(writer.BuddyOperation)(writer.BuddyErrorBody(l)(errorCode))
			}
		}
	}
}
