package compartment

import (
	"atlas-channel/kafka/message/compartment"
	"atlas-channel/kafka/producer"
	"context"
	"github.com/Chronicle20/atlas-constants/inventory"
	_map "github.com/Chronicle20/atlas-constants/map"
	"github.com/sirupsen/logrus"
)

type Processor struct {
	l   logrus.FieldLogger
	ctx context.Context
}

func NewProcessor(l logrus.FieldLogger, ctx context.Context) *Processor {
	p := &Processor{
		l:   l,
		ctx: ctx,
	}
	return p
}

func (p *Processor) Unequip(characterId uint32, inventoryType inventory.Type, source int16, destination int16) error {
	return producer.ProviderImpl(p.l)(p.ctx)(compartment.EnvCommandTopic)(UnequipAssetCommandProvider(characterId, inventoryType, source, destination))
}

func (p *Processor) Equip(characterId uint32, inventoryType inventory.Type, source int16, destination int16) error {
	return producer.ProviderImpl(p.l)(p.ctx)(compartment.EnvCommandTopic)(EquipAssetCommandProvider(characterId, inventoryType, source, destination))
}

func (p *Processor) Move(characterId uint32, inventoryType inventory.Type, source int16, destination int16) error {
	return producer.ProviderImpl(p.l)(p.ctx)(compartment.EnvCommandTopic)(MoveAssetCommandProvider(characterId, inventoryType, source, destination))
}

func (p *Processor) Drop(m _map.Model, characterId uint32, inventoryType inventory.Type, source int16, quantity int16, x int16, y int16) error {
	return producer.ProviderImpl(p.l)(p.ctx)(compartment.EnvCommandTopic)(DropAssetCommandProvider(m, characterId, inventoryType, source, quantity, x, y))
}
