package server

import (
	"fmt"
	"github.com/Chronicle20/atlas-constants/channel"
	_map "github.com/Chronicle20/atlas-constants/map"
	"github.com/Chronicle20/atlas-constants/world"
	"github.com/Chronicle20/atlas-tenant"
)

type Model struct {
	tenant    tenant.Model
	worldId   world.Id
	channelId channel.Id
}

func (m Model) Tenant() tenant.Model {
	return m.tenant
}

func (m Model) WorldId() world.Id {
	return m.worldId
}

func (m Model) ChannelId() channel.Id {
	return m.channelId
}

func (m Model) String() string {
	return fmt.Sprintf("Tenant [%s] World [%d] Channel [%d]", m.tenant.String(), m.worldId, m.channelId)
}

func New(tenant tenant.Model, worldId world.Id, channelId channel.Id) (Model, error) {
	return Model{
		tenant:    tenant,
		worldId:   worldId,
		channelId: channelId,
	}, nil
}

func (m Model) Is(t tenant.Model, worldId world.Id, channelId channel.Id) bool {
	if !m.IsWorld(t, worldId) {
		return false
	}
	if channelId != m.channelId {
		return false
	}
	return true
}

func (m Model) IsWorld(t tenant.Model, worldId world.Id) bool {
	is := t.Is(m.Tenant())
	if !is {
		return false
	}
	if worldId != m.worldId {
		return false
	}
	return true
}

func (m Model) Map(mapId _map.Id) _map.Model {
	return _map.NewModel(m.WorldId())(m.ChannelId())(mapId)
}
