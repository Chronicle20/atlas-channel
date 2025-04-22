package handler

import (
	"atlas-channel/character"
	"atlas-channel/character/skill"
	_map "atlas-channel/map"
	"atlas-channel/monster"
	"atlas-channel/session"
	skill2 "atlas-channel/skill"
	"atlas-channel/skill/effect"
	model2 "atlas-channel/socket/model"
	"atlas-channel/socket/writer"
	"context"
	"errors"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/sirupsen/logrus"
)

func processAttack(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(ai model2.AttackInfo) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(ai model2.AttackInfo) model.Operator[session.Model] {
		return func(wp writer.Producer) func(ai model2.AttackInfo) model.Operator[session.Model] {
			return func(ai model2.AttackInfo) model.Operator[session.Model] {
				return func(s session.Model) error {
					cp := character.NewProcessor(l, ctx)
					c, err := cp.GetById(cp.InventoryDecorator, cp.SkillModelDecorator)(s.CharacterId())
					if err != nil {
						return err
					}

					if ai.SkillId() > 0 {
						// Process skill
						var sk skill.Model
						for _, tsk := range c.Skills() {
							if tsk.Id() == ai.SkillId() {
								sk = tsk
							}
						}
						if sk.Id() == 0 {
							l.Errorf("Character [%d] attempting to attack with skill [%d] which they do not own.", s.CharacterId(), ai.SkillId())
							return session.NewProcessor(l, ctx).Destroy(s)
						}

						var se effect.Model
						se, err = skill2.NewProcessor(l, ctx).GetEffect(ai.SkillId(), sk.Level())
						if err != nil {
							return err
						}
						if se.HPConsume() > 0 {
							_ = cp.ChangeHP(s.Map(), s.CharacterId(), -int16(se.HPConsume()))
						}
						if se.MPConsume() > 0 {
							_ = cp.ChangeMP(s.Map(), s.CharacterId(), -int16(se.MPConsume()))
						}
					}

					for _, di := range ai.DamageInfo() {
						for _, d := range di.Damages() {
							err := monster.NewProcessor(l, ctx).Damage(s.Map(), di.MonsterId(), s.CharacterId(), d)
							if err != nil {
								l.WithError(err).Errorf("Unable to apply damage [%d] to monster [%d] from character [%d].", d, di.MonsterId(), s.CharacterId())
							}
						}
					}

					_ = _map.NewProcessor(l, ctx).ForOtherSessionsInMap(s.Map(), s.CharacterId(), func(os session.Model) error {
						var writerName string
						var bodyProducer writer.BodyProducer
						if ai.AttackType() == model2.AttackTypeMelee {
							writerName = writer.CharacterAttackMelee
							bodyProducer = writer.CharacterAttackMeleeBody(l)(ctx)(c, ai)
						} else if ai.AttackType() == model2.AttackTypeRanged {
							writerName = writer.CharacterAttackRanged
							bodyProducer = writer.CharacterAttackRangedBody(l)(ctx)(c, ai)
						} else if ai.AttackType() == model2.AttackTypeMagic {
							writerName = writer.CharacterAttackMagic
							bodyProducer = writer.CharacterAttackMagicBody(l)(ctx)(c, ai)
						} else if ai.AttackType() == model2.AttackTypeEnergy {
							writerName = writer.CharacterAttackEnergy
							bodyProducer = writer.CharacterAttackEnergyBody(l)(ctx)(c, ai)
						} else {
							return errors.New("unhandled attack type")
						}

						err = session.Announce(l)(ctx)(wp)(writerName)(bodyProducer)(os)
						if err != nil {
							l.WithError(err).Errorf("Unable to announce character [%d] is attacking to character [%d].", s.CharacterId(), os.CharacterId())
							return err
						}
						return nil
					})

					// TODO apply cooldown
					// TODO cancel dark sight / wind walk
					// TODO apply combo orbs (add or consume)
					// TODO decrease HP from DragonKnight Sacrifice
					// TODO apply attack effect (heal, mp consumption, dispel, cure all, combo reset, etc)
					// TODO destroy Chief Bandit exploded mesos
					// TODO apply Pick Pocket
					// TODO increase HP from Energy Drain, Vampire, or Drain
					// TODO apply Bandit Steal
					// TODO Fire Demon ice weaken
					// TODO Ice Demon fire weaken
					// TODO Homing Beacon / Bullseye
					// TODO Flame Thrower
					// TODO Snow Charge
					// TODO Hamstring
					// TODO Slow
					// TODO Blind
					// TODO Paladin / White Knight charges
					// TODO Combo Drain
					// TODO Mortal Blow
					// TODO Three Snails consumption
					// TODO Heavens Hammer
					// TODO ComboTempest
					// TODO BodyPressure
					// TODO Monster Weapon Atk Reflect
					// TODO Monster Magic Atk Reflect
					// TODO Apply MPEater

					return nil
				}
			}
		}
	}
}
