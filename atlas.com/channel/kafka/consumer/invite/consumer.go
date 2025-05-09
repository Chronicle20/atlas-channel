package invite

import (
	"atlas-channel/character"
	consumer2 "atlas-channel/kafka/consumer"
	invite2 "atlas-channel/kafka/message/invite"
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
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

func InitConsumers(l logrus.FieldLogger) func(func(config consumer.Config, decorators ...model.Decorator[consumer.Config])) func(consumerGroupId string) {
	return func(rf func(config consumer.Config, decorators ...model.Decorator[consumer.Config])) func(consumerGroupId string) {
		return func(consumerGroupId string) {
			rf(consumer2.NewConfig(l)("invite_status_event")(invite2.EnvEventStatusTopic)(consumerGroupId), consumer.SetHeaderParsers(consumer.SpanHeaderParser, consumer.TenantHeaderParser), consumer.SetStartOffset(kafka.LastOffset))
		}
	}
}

func InitHandlers(l logrus.FieldLogger) func(sc server.Model) func(wp writer.Producer) func(rf func(topic string, handler handler.Handler) (string, error)) {
	return func(sc server.Model) func(wp writer.Producer) func(rf func(topic string, handler handler.Handler) (string, error)) {
		return func(wp writer.Producer) func(rf func(topic string, handler handler.Handler) (string, error)) {
			return func(rf func(topic string, handler handler.Handler) (string, error)) {
				var t string
				t, _ = topic.EnvProvider(l)(invite2.EnvEventStatusTopic)()
				_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleCreatedStatusEvent(sc, wp))))
				_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleRejectedStatusEvent(sc, wp))))
			}
		}
	}
}

func handleCreatedStatusEvent(sc server.Model, wp writer.Producer) message.Handler[invite2.StatusEvent[invite2.CreatedEventBody]] {
	return func(l logrus.FieldLogger, ctx context.Context, e invite2.StatusEvent[invite2.CreatedEventBody]) {
		if e.Type != invite2.EventInviteStatusTypeCreated {
			return
		}

		if !sc.IsWorld(tenant.MustFromContext(ctx), world.Id(e.WorldId)) {
			return
		}

		rc, err := character.NewProcessor(l, ctx).GetById()(e.Body.OriginatorId)
		if err != nil {
			l.WithError(err).Errorf("Unablet to get character [%d] details, who generated the invite.", e.Body.OriginatorId)
			return
		}

		var eventHandler model.Operator[session.Model]
		if e.InviteType == invite2.InviteTypeParty {
			eventHandler = handlePartyCreatedStatusEvent(l)(ctx)(wp)(e.ReferenceId, rc.Name())
		} else if e.InviteType == invite2.InviteTypeBuddy {
			eventHandler = handleBuddyCreatedStatusEvent(l)(ctx)(wp)(e.Body.TargetId, e.ReferenceId, rc.Name())
		} else if e.InviteType == invite2.InviteTypeGuild {
			eventHandler = handleGuildCreatedStatusEvent(l)(ctx)(wp)(e.ReferenceId, rc.Name())
		} else if e.InviteType == invite2.InviteTypeMessenger {
			eventHandler = handleMessengerCreatedStatusEvent(l)(ctx)(wp)(e.ReferenceId, rc.Name())
		}

		if eventHandler != nil {
			session.NewProcessor(l, ctx).IfPresentByCharacterId(sc.WorldId(), sc.ChannelId())(e.Body.TargetId, eventHandler)
		}
	}
}

func handlePartyCreatedStatusEvent(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(partyId uint32, originatorName string) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(partyId uint32, originatorName string) model.Operator[session.Model] {
		return func(wp writer.Producer) func(partyId uint32, originatorName string) model.Operator[session.Model] {
			return func(partyId uint32, originatorName string) model.Operator[session.Model] {
				return session.Announce(l)(ctx)(wp)(writer.PartyOperation)(writer.PartyInviteBody(l)(partyId, originatorName))
			}
		}
	}
}

func handleBuddyCreatedStatusEvent(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(actorId uint32, originatorId uint32, originatorName string) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(actorId uint32, originatorId uint32, originatorName string) model.Operator[session.Model] {
		t := tenant.MustFromContext(ctx)
		return func(wp writer.Producer) func(actorId uint32, originatorId uint32, originatorName string) model.Operator[session.Model] {
			return func(actorId uint32, originatorId uint32, originatorName string) model.Operator[session.Model] {
				return session.Announce(l)(ctx)(wp)(writer.BuddyOperation)(writer.BuddyInviteBody(l, t)(actorId, originatorId, originatorName))
			}
		}
	}
}

func handleGuildCreatedStatusEvent(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(originatorId uint32, originatorName string) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(originatorId uint32, originatorName string) model.Operator[session.Model] {
		return func(wp writer.Producer) func(originatorId uint32, originatorName string) model.Operator[session.Model] {
			return func(originatorId uint32, originatorName string) model.Operator[session.Model] {
				return session.Announce(l)(ctx)(wp)(writer.GuildOperation)(writer.GuildInviteBody(l)(originatorId, originatorName))
			}
		}
	}
}

func handleMessengerCreatedStatusEvent(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(originatorId uint32, originatorName string) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(originatorId uint32, originatorName string) model.Operator[session.Model] {
		return func(wp writer.Producer) func(originatorId uint32, originatorName string) model.Operator[session.Model] {
			return func(originatorId uint32, originatorName string) model.Operator[session.Model] {
				return session.Announce(l)(ctx)(wp)(writer.MessengerOperation)(writer.MessengerOperationInviteBody(originatorName, originatorId))
			}
		}
	}
}

func handleRejectedStatusEvent(sc server.Model, wp writer.Producer) message.Handler[invite2.StatusEvent[invite2.RejectedEventBody]] {
	return func(l logrus.FieldLogger, ctx context.Context, e invite2.StatusEvent[invite2.RejectedEventBody]) {
		if e.Type != invite2.EventInviteStatusTypeRejected {
			return
		}

		if !sc.IsWorld(tenant.MustFromContext(ctx), world.Id(e.WorldId)) {
			return
		}

		rc, err := character.NewProcessor(l, ctx).GetById()(e.Body.TargetId)
		if err != nil {
			l.WithError(err).Errorf("Unablet to get character [%d] details, who generated the invite.", e.Body.OriginatorId)
			return
		}

		var eventHandler model.Operator[session.Model]
		if e.InviteType == invite2.InviteTypeParty {
			eventHandler = handlePartyRejectedStatusEvent(l)(ctx)(wp)(rc.Name())
		} else if e.InviteType == invite2.InviteTypeBuddy {
			// TODO send rejection to requesting character.
		} else if e.InviteType == invite2.InviteTypeGuild {
			eventHandler = handleGuildRejectedStatusEvent(l)(ctx)(wp)(rc.Name())
		} else if e.InviteType == invite2.InviteTypeMessenger {
			eventHandler = handleMessengerRejectedStatusEvent(l)(ctx)(wp)(rc.Name())
		}

		if eventHandler != nil {
			session.NewProcessor(l, ctx).IfPresentByCharacterId(sc.WorldId(), sc.ChannelId())(e.Body.OriginatorId, eventHandler)
		}
	}
}

func handlePartyRejectedStatusEvent(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(targetName string) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(targetName string) model.Operator[session.Model] {
		return func(wp writer.Producer) func(targetName string) model.Operator[session.Model] {
			return func(targetName string) model.Operator[session.Model] {
				return session.Announce(l)(ctx)(wp)(writer.PartyOperation)(writer.PartyErrorBody(l)("HAVE_DENIED_REQUEST_TO_THE_PARTY", targetName))
			}
		}
	}
}

func handleGuildRejectedStatusEvent(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(targetName string) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(targetName string) model.Operator[session.Model] {
		return func(wp writer.Producer) func(targetName string) model.Operator[session.Model] {
			return func(targetName string) model.Operator[session.Model] {
				return session.Announce(l)(ctx)(wp)(writer.GuildOperation)(writer.GuildErrorBody2(l)(writer.GuildOperationInviteDenied, targetName))
			}
		}
	}
}

func handleMessengerRejectedStatusEvent(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(targetName string) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(targetName string) model.Operator[session.Model] {
		return func(wp writer.Producer) func(targetName string) model.Operator[session.Model] {
			return func(targetName string) model.Operator[session.Model] {
				return session.Announce(l)(ctx)(wp)(writer.MessengerOperation)(writer.MessengerOperationInviteDeclinedBody(targetName, 0))
			}
		}
	}
}
