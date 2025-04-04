package pet

import (
	"atlas-channel/kafka/producer"
	"atlas-channel/pet/exclude"
	"context"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/Chronicle20/atlas-rest/requests"
	"github.com/sirupsen/logrus"
)

func ByIdProvider(l logrus.FieldLogger) func(ctx context.Context) func(petId uint64) model.Provider[Model] {
	return func(ctx context.Context) func(petId uint64) model.Provider[Model] {
		return func(petId uint64) model.Provider[Model] {
			return requests.Provider[RestModel, Model](l, ctx)(requestById(petId), Extract)
		}
	}
}

func GetById(l logrus.FieldLogger) func(ctx context.Context) func(petId uint64) (Model, error) {
	return func(ctx context.Context) func(petId uint64) (Model, error) {
		return func(petId uint64) (Model, error) {
			return ByIdProvider(l)(ctx)(petId)()
		}
	}
}

func ByOwnerProvider(l logrus.FieldLogger) func(ctx context.Context) func(ownerId uint32) model.Provider[[]Model] {
	return func(ctx context.Context) func(ownerId uint32) model.Provider[[]Model] {
		return func(ownerId uint32) model.Provider[[]Model] {
			return requests.SliceProvider[RestModel, Model](l, ctx)(requestByOwnerId(ownerId), Extract, model.Filters[Model]())
		}
	}
}

func GetByOwner(l logrus.FieldLogger) func(ctx context.Context) func(ownerId uint32) ([]Model, error) {
	return func(ctx context.Context) func(ownerId uint32) ([]Model, error) {
		return func(ownerId uint32) ([]Model, error) {
			return ByOwnerProvider(l)(ctx)(ownerId)()
		}
	}
}

func GetByOwnerItem(l logrus.FieldLogger) func(ctx context.Context) func(ownerId uint32, inventoryItemId uint32) (Model, error) {
	return func(ctx context.Context) func(ownerId uint32, inventoryItemId uint32) (Model, error) {
		return func(ownerId uint32, inventoryItemId uint32) (Model, error) {
			return model.FirstProvider(ByOwnerProvider(l)(ctx)(ownerId), model.Filters(InventoryItemIdFilter(inventoryItemId)))()
		}
	}
}

func InventoryItemIdFilter(inventoryItemId uint32) model.Filter[Model] {
	return func(m Model) bool {
		return m.InventoryItemId() == inventoryItemId
	}
}

func Spawn(l logrus.FieldLogger) func(ctx context.Context) func(characterId uint32, petId uint64, lead bool) error {
	return func(ctx context.Context) func(characterId uint32, petId uint64, lead bool) error {
		return func(characterId uint32, petId uint64, lead bool) error {
			l.Debugf("Character [%d] attempting to spawn pet [%d]", characterId, petId)
			return producer.ProviderImpl(l)(ctx)(EnvCommandTopic)(spawnProvider(characterId, petId, lead))
		}
	}
}

func Despawn(l logrus.FieldLogger) func(ctx context.Context) func(characterId uint32, petId uint64) error {
	return func(ctx context.Context) func(characterId uint32, petId uint64) error {
		return func(characterId uint32, petId uint64) error {
			l.Debugf("Character [%d] attempting to despawn pet [%d].", characterId, petId)
			return producer.ProviderImpl(l)(ctx)(EnvCommandTopic)(despawnProvider(characterId, petId))
		}
	}
}

func AttemptCommand(l logrus.FieldLogger) func(ctx context.Context) func(petId uint64, commandId byte, byName bool, characterId uint32) error {
	return func(ctx context.Context) func(petId uint64, commandId byte, byName bool, characterId uint32) error {
		return func(petId uint64, commandId byte, byName bool, characterId uint32) error {
			l.Debugf("Character [%d] triggered pet [%d] command. byName [%t], command [%d]", characterId, petId, byName, commandId)
			return producer.ProviderImpl(l)(ctx)(EnvCommandTopic)(attemptCommandProvider(petId, commandId, byName, characterId))
		}
	}
}

func SetExcludeItems(l logrus.FieldLogger) func(ctx context.Context) func(characterId uint32, petId uint64, items []exclude.Model) error {
	return func(ctx context.Context) func(characterId uint32, petId uint64, items []exclude.Model) error {
		return func(characterId uint32, petId uint64, items []exclude.Model) error {
			l.Debugf("Character [%d] setting exclude items for pet [%d]. count [%d].", characterId, petId, len(items))
			itemIds := make([]uint32, len(items))
			for i, item := range items {
				itemIds[i] = item.ItemId()
			}
			return producer.ProviderImpl(l)(ctx)(EnvCommandTopic)(setExcludesCommandProvider(characterId, petId, itemIds))
		}
	}
}
