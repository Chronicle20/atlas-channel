package session

import (
	"atlas-channel/buddylist"
	"atlas-channel/character"
	"atlas-channel/character/buff"
	"atlas-channel/character/key"
	"atlas-channel/guild"
	consumer2 "atlas-channel/kafka/consumer"
	session2 "atlas-channel/kafka/message/account/session"
	"atlas-channel/macro"
	"atlas-channel/server"
	"atlas-channel/session"
	model2 "atlas-channel/socket/model"
	"atlas-channel/socket/writer"
	"atlas-channel/world"
	"context"
	_map "github.com/Chronicle20/atlas-constants/map"
	"github.com/Chronicle20/atlas-kafka/consumer"
	"github.com/Chronicle20/atlas-kafka/handler"
	"github.com/Chronicle20/atlas-kafka/message"
	"github.com/Chronicle20/atlas-kafka/topic"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/Chronicle20/atlas-tenant"
	"github.com/sirupsen/logrus"
	"sort"
)

func InitConsumers(l logrus.FieldLogger) func(func(config consumer.Config, decorators ...model.Decorator[consumer.Config])) func(consumerGroupId string) {
	return func(rf func(config consumer.Config, decorators ...model.Decorator[consumer.Config])) func(consumerGroupId string) {
		return func(consumerGroupId string) {
			rf(consumer2.NewConfig(l)("account_session_status_event")(session2.EnvEventStatusTopic)(consumerGroupId), consumer.SetHeaderParsers(consumer.SpanHeaderParser, consumer.TenantHeaderParser))
		}
	}
}

func InitHandlers(l logrus.FieldLogger) func(sc server.Model) func(wp writer.Producer) func(rf func(topic string, handler handler.Handler) (string, error)) {
	return func(sc server.Model) func(wp writer.Producer) func(rf func(topic string, handler handler.Handler) (string, error)) {
		return func(wp writer.Producer) func(rf func(topic string, handler handler.Handler) (string, error)) {
			return func(rf func(topic string, handler handler.Handler) (string, error)) {
				var t string
				t, _ = topic.EnvProvider(l)(session2.EnvEventStatusTopic)()
				_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleError(sc, wp))))
				_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleChannelChange(sc, wp))))
				_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handlePlayerLoggedIn(sc, wp))))
			}
		}
	}
}

func handleError(sc server.Model, wp writer.Producer) func(l logrus.FieldLogger, ctx context.Context, e session2.StatusEvent[session2.ErrorStatusEventBody]) {
	return func(l logrus.FieldLogger, ctx context.Context, e session2.StatusEvent[session2.ErrorStatusEventBody]) {
		if e.Type != session2.EventStatusTypeError {
			return
		}

		t := sc.Tenant()
		if !t.Is(tenant.MustFromContext(ctx)) {
			return
		}

		session.NewProcessor(l, ctx).IfPresentByIdInWorld(e.SessionId, sc.WorldId(), sc.ChannelId(), announceError(l)(ctx)(wp)(e.Body.Code))
	}
}

func announceError(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(reason string) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(reason string) model.Operator[session.Model] {
		return func(wp writer.Producer) func(reason string) model.Operator[session.Model] {
			return func(reason string) model.Operator[session.Model] {
				return func(s session.Model) error {
					l.Errorf("Unable to update session for character [%d] attempting to switch to channel.", s.CharacterId())
					return session.NewProcessor(l, ctx).Destroy(s)
				}
			}
		}
	}
}

func handleChannelChange(sc server.Model, wp writer.Producer) message.Handler[session2.StatusEvent[session2.StateChangedEventBody[model2.ChannelChange]]] {
	return func(l logrus.FieldLogger, ctx context.Context, e session2.StatusEvent[session2.StateChangedEventBody[model2.ChannelChange]]) {
		if e.Type != session2.EventStatusTypeStateChanged {
			return
		}

		t := sc.Tenant()
		if !t.Is(tenant.MustFromContext(ctx)) {
			return
		}

		if len(e.Body.Params.IPAddress) <= 0 {
			return
		}

		session.NewProcessor(l, ctx).IfPresentByIdInWorld(e.SessionId, sc.WorldId(), sc.ChannelId(), processChannelChangeReturn(l)(ctx)(wp)(e.AccountId, e.Body.State, e.Body.Params))
	}
}

func processChannelChangeReturn(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(accountId uint32, state uint8, params model2.ChannelChange) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(accountId uint32, state uint8, params model2.ChannelChange) model.Operator[session.Model] {
		return func(wp writer.Producer) func(accountId uint32, state uint8, params model2.ChannelChange) model.Operator[session.Model] {
			return func(accountId uint32, state uint8, params model2.ChannelChange) model.Operator[session.Model] {
				return session.Announce(l)(ctx)(wp)(writer.ChannelChange)(writer.ChannelChangeBody(params.IPAddress, params.Port))
			}
		}
	}
}

func handlePlayerLoggedIn(sc server.Model, wp writer.Producer) message.Handler[session2.StatusEvent[session2.StateChangedEventBody[model2.SetField]]] {
	return func(l logrus.FieldLogger, ctx context.Context, e session2.StatusEvent[session2.StateChangedEventBody[model2.SetField]]) {
		if e.Type != session2.EventStatusTypeStateChanged {
			return
		}

		t := sc.Tenant()
		if !t.Is(tenant.MustFromContext(ctx)) {
			return
		}

		session.NewProcessor(l, ctx).IfPresentByIdInWorld(e.SessionId, sc.WorldId(), sc.ChannelId(), processStateReturn(l)(ctx)(wp)(e.AccountId, e.Body.State, e.Body.Params))
	}
}

func processStateReturn(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(accountId uint32, state uint8, params model2.SetField) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(accountId uint32, state uint8, params model2.SetField) model.Operator[session.Model] {
		sp := session.NewProcessor(l, ctx)
		t := tenant.MustFromContext(ctx)
		return func(wp writer.Producer) func(accountId uint32, state uint8, params model2.SetField) model.Operator[session.Model] {
			return func(accountId uint32, state uint8, params model2.SetField) model.Operator[session.Model] {
				return func(s session.Model) error {
					if params.CharacterId <= 0 {
						return nil
					}

					cp := character.NewProcessor(l, ctx)
					c, err := cp.GetById(cp.InventoryDecorator, cp.SkillModelDecorator, cp.PetModelDecorator)(params.CharacterId)
					if err != nil {
						l.WithError(err).Errorf("Unable to locate character [%d] attempting to login.", params.CharacterId)
						return sp.Destroy(s)
					}
					bl, err := buddylist.NewProcessor(l, ctx).GetById(params.CharacterId)
					if err != nil {
						l.WithError(err).Errorf("Unable to locate buddylist [%d] attempting to login.", params.CharacterId)
						return sp.Destroy(s)
					}

					s = sp.SetAccountId(s.SessionId(), c.AccountId())
					s = sp.SetCharacterId(s.SessionId(), c.Id())
					s = sp.SetGm(s.SessionId(), c.Gm())
					s = sp.SetMapId(s.SessionId(), _map.Id(c.MapId()))

					sp.SessionCreated(s)

					l.Debugf("Writing SetField for character [%d].", c.Id())
					err = session.Announce(l)(ctx)(wp)(writer.SetField)(writer.SetFieldBody(l, t)(s.ChannelId(), c, bl))(s)
					if err != nil {
						l.WithError(err).Errorf("Unable to show set field response for character [%d]", c.Id())
					}
					go func() {
						err := session.Announce(l)(ctx)(wp)(writer.BuddyOperation)(writer.BuddyListUpdateBody(l, t)(bl.Buddies()))(s)
						if err != nil {
							l.WithError(err).Errorf("Unable to write character [%d] buddy list.", c.Id())
						}
					}()
					go func() {
						g, _ := guild.NewProcessor(l, ctx).GetByMemberId(c.Id())
						if g.Id() != 0 {
							err := session.Announce(l)(ctx)(wp)(writer.GuildOperation)(writer.GuildInfoBody(l, t)(g))(s)
							if err != nil {
								l.WithError(err).Errorf("Unable to write character [%d] buddy list.", c.Id())
							}
						}
					}()
					go func() {
						km, err := model.CollectToMap[key.Model, int32, key.Model](key.NewProcessor(l, ctx).ByCharacterIdProvider(s.CharacterId()), func(m key.Model) int32 {
							return m.Key()
						}, func(m key.Model) key.Model {
							return m
						})()
						if err != nil {
							l.WithError(err).Errorf("Unable to show key map for character [%d].", s.CharacterId())
							return
						}

						err = session.Announce(l)(ctx)(wp)(writer.CharacterKeyMap)(writer.CharacterKeyMapBody(km))(s)
						if err != nil {
							l.WithError(err).Errorf("Unable to show key map for character [%d].", s.CharacterId())
						}
					}()
					go func() {
						bs, err := buff.NewProcessor(l, ctx).GetByCharacterId(s.CharacterId())
						if err != nil {
							l.WithError(err).Errorf("Unable to retrieve active buffs for character [%d].", s.CharacterId())
							return
						}
						err = session.Announce(l)(ctx)(wp)(writer.CharacterBuffGive)(writer.CharacterBuffGiveBody(l)(ctx)(bs))(s)
						if err != nil {
							l.WithError(err).Errorf("Unable to write character [%d] buddy list.", c.Id())
						}
					}()
					go func() {
						w, err := world.NewProcessor(l, ctx).GetById(byte(s.WorldId()))
						if err != nil {
							return
						}
						err = session.Announce(l)(ctx)(wp)(writer.WorldMessage)(writer.WorldMessageTopScrollBody(l, t)(w.Message()))(s)
						if err != nil {
							l.WithError(err).Errorf("Unable to write character [%d] buddy list.", c.Id())
						}
					}()
					go func() {
						var sms []macro.Model
						sms, err = macro.NewProcessor(l, ctx).GetByCharacterId(s.CharacterId())
						if err != nil {
							l.WithError(err).Errorf("Unable to read skill macros for character [%d].", c.Id())
							return
						}
						sort.Slice(sms, func(i, j int) bool {
							return sms[i].Id() < sms[j].Id()
						})
						mms := make([]model2.Macro, 0)
						for _, sm := range sms {
							mms = append(mms, model2.NewMacro(sm.Name(), sm.Shout(), sm.SkillId1(), sm.SkillId2(), sm.SkillId3()))
						}
						err = session.Announce(l)(ctx)(wp)(writer.CharacterSkillMacro)(writer.CharacterSkillMacroBody(l, tenant.MustFromContext(ctx))(model2.NewMacros(mms...)))(s)
						if err != nil {
							l.WithError(err).Errorf("Unable to show key map for character [%d].", s.CharacterId())
						}
					}()
					return nil
				}
			}
		}
	}
}
