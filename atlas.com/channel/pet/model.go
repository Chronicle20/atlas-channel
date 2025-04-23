package pet

import (
	"atlas-channel/pet/exclude"
	"time"
)

type Model struct {
	id         uint32
	cashId     uint64
	templateId uint32
	name       string
	level      byte
	closeness  uint16
	fullness   byte
	expiration time.Time
	ownerId    uint32
	slot       int8
	x          int16
	y          int16
	stance     byte
	fh         int16
	excludes   []exclude.Model
	flag       uint16
	purchaseBy uint32
}

func (m Model) Id() uint32 {
	return m.id
}

func (m Model) CashId() uint64 {
	return m.cashId
}

func (m Model) TemplateId() uint32 {
	return m.templateId
}

func (m Model) Name() string {
	return m.name
}

func (m Model) Level() byte {
	return m.level
}

func (m Model) Closeness() uint16 {
	return m.closeness
}

func (m Model) Fullness() byte {
	return m.fullness
}

func (m Model) Expiration() time.Time {
	return m.expiration
}

func (m Model) OwnerId() uint32 {
	return m.ownerId
}

func (m Model) Lead() bool {
	return m.slot == 0
}

func (m Model) Slot() int8 {
	return m.slot
}

func (m Model) X() int16 {
	return m.x
}

func (m Model) Y() int16 {
	return m.y
}

func (m Model) Stance() byte {
	return m.stance
}

func (m Model) Fh() int16 {
	return m.fh
}

func (m Model) Excludes() []exclude.Model {
	return m.excludes
}

func (m Model) Flag() uint16 {
	return m.flag
}

func (m Model) PurchaseBy() uint32 {
	return m.purchaseBy
}

type ModelBuilder struct {
	id         uint32
	cashId     uint64
	templateId uint32
	name       string
	level      byte
	closeness  uint16
	fullness   byte
	expiration time.Time
	ownerId    uint32
	slot       int8
	x          int16
	y          int16
	stance     byte
	fh         int16
	excludes   []exclude.Model
	flag       uint16
	purchaseBy uint32
}

func NewModelBuilder(id uint32, cashId uint64, templateId uint32, name string) *ModelBuilder {
	return &ModelBuilder{
		id:         id,
		cashId:     cashId,
		templateId: templateId,
		name:       name,
		excludes:   make([]exclude.Model, 0),
	}
}

func (b *ModelBuilder) SetLevel(level byte) *ModelBuilder {
	b.level = level
	return b
}

func (b *ModelBuilder) SetCloseness(closeness uint16) *ModelBuilder {
	b.closeness = closeness
	return b
}

func (b *ModelBuilder) SetFullness(fullness byte) *ModelBuilder {
	b.fullness = fullness
	return b
}

func (b *ModelBuilder) SetExpiration(expiration time.Time) *ModelBuilder {
	b.expiration = expiration
	return b
}

func (b *ModelBuilder) SetOwnerID(ownerId uint32) *ModelBuilder {
	b.ownerId = ownerId
	return b
}

func (b *ModelBuilder) SetSlot(slot int8) *ModelBuilder {
	b.slot = slot
	return b
}

func (b *ModelBuilder) SetX(x int16) *ModelBuilder {
	b.x = x
	return b
}

func (b *ModelBuilder) SetY(y int16) *ModelBuilder {
	b.y = y
	return b
}

func (b *ModelBuilder) SetStance(stance byte) *ModelBuilder {
	b.stance = stance
	return b
}

func (b *ModelBuilder) SetExcludes(excludes []exclude.Model) *ModelBuilder {
	b.excludes = excludes
	return b
}

func (b *ModelBuilder) SetFoothold(fh int16) *ModelBuilder {
	b.fh = fh
	return b
}

func (b *ModelBuilder) Build() Model {
	return Model{
		id:         b.id,
		cashId:     b.cashId,
		templateId: b.templateId,
		name:       b.name,
		level:      b.level,
		closeness:  b.closeness,
		fullness:   b.fullness,
		expiration: b.expiration,
		ownerId:    b.ownerId,
		slot:       b.slot,
		x:          b.x,
		y:          b.y,
		stance:     b.stance,
		fh:         b.fh,
		excludes:   b.excludes,
		flag:       b.flag,
		purchaseBy: b.purchaseBy,
	}
}
