package commodities

import "github.com/google/uuid"

type Model struct {
	id           uuid.UUID
	templateId   uint32
	mesoPrice    uint32
	discountRate byte
	tokenItemId  uint32
	tokenPrice   uint32
	period       uint32
	levelLimit   uint32
}

// Id returns the model's id
func (m *Model) Id() uuid.UUID {
	return m.id
}

// TemplateId returns the model's templateId
func (m *Model) TemplateId() uint32 {
	return m.templateId
}

// MesoPrice returns the model's mesoPrice
func (m *Model) MesoPrice() uint32 {
	return m.mesoPrice
}

// DiscountRate returns the model's discountRate
func (m *Model) DiscountRate() byte {
	return m.discountRate
}

// TokenItemId returns the model's tokenItemId
func (m *Model) TokenItemId() uint32 {
	return m.tokenItemId
}

// TokenPrice returns the model's tokenPrice
func (m *Model) TokenPrice() uint32 {
	return m.tokenPrice
}

// Period returns the model's period
func (m *Model) Period() uint32 {
	return m.period
}

// LevelLimit returns the model's levelLimit
func (m *Model) LevelLimit() uint32 {
	return m.levelLimit
}

func (m *Model) UnitPrice() uint64 {
	// TODO
	return 0
}

func (m *Model) Quantity() uint16 {
	// TODO
	return 0
}

func (m *Model) MaxPerSlot() uint16 {
	// TODO
	return 0
}

// ModelBuilder is used to build Model instances
type ModelBuilder struct {
	id           uuid.UUID
	templateId   uint32
	mesoPrice    uint32
	discountRate byte
	tokenItemId  uint32
	tokenPrice   uint32
	period       uint32
	levelLimit   uint32
}

// SetId sets the id for the ModelBuilder
func (b *ModelBuilder) SetId(id uuid.UUID) *ModelBuilder {
	b.id = id
	return b
}

// SetTemplateId sets the templateId for the ModelBuilder
func (b *ModelBuilder) SetTemplateId(templateId uint32) *ModelBuilder {
	b.templateId = templateId
	return b
}

// SetMesoPrice sets the mesoPrice for the ModelBuilder
func (b *ModelBuilder) SetMesoPrice(mesoPrice uint32) *ModelBuilder {
	b.mesoPrice = mesoPrice
	return b
}

// SetDiscountRate sets the discountRate for the ModelBuilder
func (b *ModelBuilder) SetDiscountRate(discountRate byte) *ModelBuilder {
	b.discountRate = discountRate
	return b
}

// SetTokenItemId sets the tokenItemId for the ModelBuilder
func (b *ModelBuilder) SetTokenItemId(tokenItemId uint32) *ModelBuilder {
	b.tokenItemId = tokenItemId
	return b
}

// SetTokenPrice sets the tokenPrice for the ModelBuilder
func (b *ModelBuilder) SetTokenPrice(tokenPrice uint32) *ModelBuilder {
	b.tokenPrice = tokenPrice
	return b
}

// SetPeriod sets the period for the ModelBuilder
func (b *ModelBuilder) SetPeriod(period uint32) *ModelBuilder {
	b.period = period
	return b
}

// SetLevelLimit sets the levelLimit for the ModelBuilder
func (b *ModelBuilder) SetLevelLimit(levelLimited uint32) *ModelBuilder {
	b.levelLimit = levelLimited
	return b
}

// Build creates a new Model instance with the builder's values
func (b *ModelBuilder) Build() Model {
	return Model{
		id:           b.id,
		templateId:   b.templateId,
		mesoPrice:    b.mesoPrice,
		discountRate: b.discountRate,
		tokenItemId:  b.tokenItemId,
		tokenPrice:   b.tokenPrice,
		period:       b.period,
		levelLimit:   b.levelLimit,
	}
}

// Clone creates a new ModelBuilder with values from the given Model
func Clone(m Model) *ModelBuilder {
	return &ModelBuilder{
		id:           m.id,
		templateId:   m.templateId,
		mesoPrice:    m.mesoPrice,
		discountRate: m.discountRate,
		tokenItemId:  m.tokenItemId,
		tokenPrice:   m.tokenPrice,
		period:       m.period,
		levelLimit:   m.levelLimit,
	}
}
