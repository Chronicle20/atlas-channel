package monster

import (
	"github.com/Chronicle20/atlas-constants/channel"
	_map "github.com/Chronicle20/atlas-constants/map"
	"github.com/Chronicle20/atlas-constants/world"
)

type Model struct {
	worldId            world.Id
	channelId          channel.Id
	mapId              _map.Id
	uniqueId           uint32
	maxHp              uint32
	hp                 uint32
	mp                 uint32
	monsterId          uint32
	controlCharacterId uint32
	x                  int16
	y                  int16
	fh                 int16
	stance             byte
	team               int8
}

func (m Model) UniqueId() uint32 {
	return m.uniqueId
}

func (m Model) Controlled() bool {
	return m.controlCharacterId != 0
}

func (m Model) MonsterId() uint32 {
	return m.monsterId
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

func (m Model) Team() int8 {
	return m.team
}

func (m Model) WorldId() world.Id {
	return m.worldId
}

func (m Model) ChannelId() channel.Id {
	return m.channelId
}

func (m Model) MapId() _map.Id {
	return m.mapId
}

func (m Model) MP() uint32 {
	return m.mp
}

func (m Model) HP() uint32 {
	return m.hp
}

func (m Model) MaxHP() uint32 {
	return m.maxHp
}
