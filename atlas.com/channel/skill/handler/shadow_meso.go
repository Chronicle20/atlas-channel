package handler

import (
	"atlas-channel/character/buff"
	"atlas-channel/skill/effect"
	"atlas-channel/socket/model"
	"context"
	_map "github.com/Chronicle20/atlas-constants/map"
	"github.com/sirupsen/logrus"
)

func UseShadowMeso(l logrus.FieldLogger) func(ctx context.Context) func(m _map.Model, characterId uint32, info model.SkillUsageInfo, effect effect.Model) error {
	return func(ctx context.Context) func(m _map.Model, characterId uint32, info model.SkillUsageInfo, effect effect.Model) error {
		return func(m _map.Model, characterId uint32, info model.SkillUsageInfo, effect effect.Model) error {

			_ = buff.Apply(l)(ctx)(m, characterId, int32(info.SkillId()), effect.Duration(), effect.StatUps())(characterId)

			return nil
		}
	}
}
