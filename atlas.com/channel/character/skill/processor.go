package skill

import (
	"atlas-channel/kafka/producer"
	"context"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/Chronicle20/atlas-rest/requests"
	"github.com/sirupsen/logrus"
)

func byCharacterIdProvider(l logrus.FieldLogger) func(ctx context.Context) func(characterId uint32) model.Provider[[]Model] {
	return func(ctx context.Context) func(characterId uint32) model.Provider[[]Model] {
		return func(characterId uint32) model.Provider[[]Model] {
			return requests.SliceProvider[RestModel, Model](l, ctx)(requestByCharacterId(characterId), Extract, model.Filters[Model]())
		}
	}
}

func GetByCharacterId(l logrus.FieldLogger) func(ctx context.Context) func(characterId uint32) ([]Model, error) {
	return func(ctx context.Context) func(characterId uint32) ([]Model, error) {
		return func(characterId uint32) ([]Model, error) {
			return byCharacterIdProvider(l)(ctx)(characterId)()
		}
	}
}

func ApplyCooldown(l logrus.FieldLogger) func(ctx context.Context) func(worldId byte, channelId byte, skillId uint32, cooldown uint32) model.Operator[uint32] {
	return func(ctx context.Context) func(worldId byte, channelId byte, skillId uint32, cooldown uint32) model.Operator[uint32] {
		return func(worldId byte, channelId byte, skillId uint32, cooldown uint32) model.Operator[uint32] {
			return func(characterId uint32) error {
				return producer.ProviderImpl(l)(ctx)(EnvCommandTopic)(setCooldownCommandProvider(characterId, skillId, cooldown))
			}
		}
	}
}
