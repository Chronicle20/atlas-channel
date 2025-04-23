package character

import (
	"atlas-channel/asset"
	"atlas-channel/character/skill"
	"atlas-channel/inventory"
	character2 "atlas-channel/kafka/message/character"
	"atlas-channel/kafka/producer"
	"atlas-channel/pet"
	"context"
	"errors"
	inventory2 "github.com/Chronicle20/atlas-constants/inventory"
	_map "github.com/Chronicle20/atlas-constants/map"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/Chronicle20/atlas-rest/requests"
	"github.com/sirupsen/logrus"
	"sort"
)

type Processor struct {
	l   logrus.FieldLogger
	ctx context.Context
	ip  *inventory.Processor
}

func NewProcessor(l logrus.FieldLogger, ctx context.Context) *Processor {
	p := &Processor{
		l:   l,
		ctx: ctx,
		ip:  inventory.NewProcessor(l, ctx),
	}
	return p
}

func (p *Processor) GetById(decorators ...model.Decorator[Model]) func(characterId uint32) (Model, error) {
	return func(characterId uint32) (Model, error) {
		mp := requests.Provider[RestModel, Model](p.l, p.ctx)(requestById(characterId), Extract)
		return model.Map(model.Decorate(decorators))(mp)()
	}
}

// deprecated
func (p *Processor) PetModelDecorator(m Model) Model {
	ms, err := pet.NewProcessor(p.l, p.ctx).GetByOwner(m.Id())
	if err != nil {
		return m
	}
	if len(ms) == 0 {
		return m
	}
	sort.Slice(ms, func(i, j int) bool {
		return ms[i].Slot() < ms[j].Slot()
	})
	return m.SetPets(ms)
}

func (p *Processor) InventoryDecorator(m Model) Model {
	i, err := p.ip.GetByCharacterId(m.Id())
	if err != nil {
		return m
	}
	return m.SetInventory(i)
}

func (p *Processor) SkillModelDecorator(m Model) Model {
	ms, err := skill.NewProcessor(p.l, p.ctx).GetByCharacterId(m.Id())
	if err != nil {
		return m
	}
	return m.SetSkills(ms)
}

func (p *Processor) GetEquipableInSlot(characterId uint32, slot int16) model.Provider[asset.Model[any]] {
	// TODO this needs to be more performant
	c, err := p.GetById(p.InventoryDecorator)(characterId)
	if err != nil {
		return model.ErrorProvider[asset.Model[any]](err)
	}
	for _, e := range c.Inventory().Equipable().Assets() {
		if e.Slot() == slot {
			return model.FixedProvider(e)
		}
	}
	return model.ErrorProvider[asset.Model[any]](errors.New("equipable not found"))
}

func (p *Processor) GetItemInSlot(characterId uint32, inventoryType inventory2.Type, slot int16) model.Provider[asset.Model[any]] {
	// TODO this needs to be more performant
	c, err := p.GetById(p.InventoryDecorator)(characterId)
	if err != nil {
		return model.ErrorProvider[asset.Model[any]](err)
	}

	cm := c.Inventory().CompartmentByType(inventoryType)
	for _, e := range cm.Assets() {
		if e.Slot() == slot {
			return model.FixedProvider(e)
		}
	}
	return model.ErrorProvider[asset.Model[any]](errors.New("item not found"))
}

func (p *Processor) ByNameProvider(name string) model.Provider[[]Model] {
	return requests.SliceProvider[RestModel, Model](p.l, p.ctx)(requestByName(name), Extract, model.Filters[Model]())
}

func (p *Processor) GetByName(name string) (Model, error) {
	return model.FirstProvider(p.ByNameProvider(name), model.Filters[Model]())()
}

type DistributePacket struct {
	Flag  uint32
	Value uint32
}

func (p *Processor) RequestDistributeAp(m _map.Model, characterId uint32, updateTime uint32, distributes []DistributePacket) error {
	var distributions = make([]character2.DistributePair, 0)
	for _, d := range distributes {
		a, err := abilityFromFlag(d.Flag)
		if err != nil {
			p.l.WithError(err).Errorf("Character [%d] passed invalid flag when attempting to distribute AP.", characterId)
			return err
		}

		distributions = append(distributions, character2.DistributePair{
			Ability: a,
			Amount:  int8(d.Value),
		})
	}
	return producer.ProviderImpl(p.l)(p.ctx)(character2.EnvCommandTopic)(RequestDistributeApCommandProvider(m, characterId, distributions))
}

func abilityFromFlag(flag uint32) (string, error) {
	switch flag {
	case 64:
		return character2.CommandDistributeApAbilityStrength, nil
	case 128:
		return character2.CommandDistributeApAbilityDexterity, nil
	case 256:
		return character2.CommandDistributeApAbilityIntelligence, nil
	case 512:
		return character2.CommandDistributeApAbilityLuck, nil
	case 2048:
		return character2.CommandDistributeApAbilityHp, nil
	case 8192:
		return character2.CommandDistributeApAbilityMp, nil
	}
	return "", errors.New("invalid flag")
}

func (p *Processor) RequestDropMeso(m _map.Model, characterId uint32, amount uint32) error {
	return producer.ProviderImpl(p.l)(p.ctx)(character2.EnvCommandTopic)(RequestDropMesoCommandProvider(m, characterId, amount))
}

func (p *Processor) ChangeHP(m _map.Model, characterId uint32, amount int16) error {
	return producer.ProviderImpl(p.l)(p.ctx)(character2.EnvCommandTopic)(ChangeHPCommandProvider(m, characterId, amount))
}

func (p *Processor) ChangeMP(m _map.Model, characterId uint32, amount int16) error {
	return producer.ProviderImpl(p.l)(p.ctx)(character2.EnvCommandTopic)(ChangeMPCommandProvider(m, characterId, amount))
}

func (p *Processor) RequestDistributeSp(m _map.Model, characterId uint32, updateTime uint32, skillId uint32, amount int8) error {
	return producer.ProviderImpl(p.l)(p.ctx)(character2.EnvCommandTopic)(RequestDistributeSpCommandProvider(m, characterId, skillId, amount))
}
