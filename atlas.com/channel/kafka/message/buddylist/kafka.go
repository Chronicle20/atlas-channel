package buddylist

import (
	"github.com/Chronicle20/atlas-constants/world"
)

const (
	EnvCommandTopic          = "COMMAND_TOPIC_BUDDY_LIST"
	CommandTypeRequestAdd    = "REQUEST_ADD"
	CommandTypeRequestDelete = "REQUEST_DELETE"
)

type Command[E any] struct {
	WorldId     world.Id `json:"worldId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
	Body        E      `json:"body"`
}

type RequestAddBuddyCommandBody struct {
	CharacterId uint32 `json:"characterId"`
	Group       string `json:"group"`
}

type RequestDeleteBuddyCommandBody struct {
	CharacterId uint32 `json:"characterId"`
}

const (
	EnvStatusEventTopic                = "EVENT_TOPIC_BUDDY_LIST_STATUS"
	StatusEventTypeBuddyAdded          = "BUDDY_ADDED"
	StatusEventTypeBuddyRemoved        = "BUDDY_REMOVED"
	StatusEventTypeBuddyUpdated        = "BUDDY_UPDATED"
	StatusEventTypeBuddyChannelChange  = "BUDDY_CHANNEL_CHANGE"
	StatusEventTypeBuddyCapacityUpdate = "CAPACITY_CHANGE"
	StatusEventTypeError               = "ERROR"

	StatusEventErrorListFull          = "BUDDY_LIST_FULL"
	StatusEventErrorOtherListFull     = "OTHER_BUDDY_LIST_FULL"
	StatusEventErrorAlreadyBuddy      = "ALREADY_BUDDY"
	StatusEventErrorCannotBuddyGm     = "CANNOT_BUDDY_GM"
	StatusEventErrorCharacterNotFound = "CHARACTER_NOT_FOUND"
	StatusEventErrorUnknownError      = "UNKNOWN_ERROR"
)

type StatusEvent[E any] struct {
	WorldId     world.Id `json:"worldId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
	Body        E      `json:"body"`
}

type BuddyAddedStatusEventBody struct {
	CharacterId   uint32 `json:"characterId"`
	Group         string `json:"group"`
	CharacterName string `json:"characterName"`
	ChannelId     int8   `json:"channelId"`
}

type BuddyRemovedStatusEventBody struct {
	CharacterId uint32 `json:"characterId"`
}

type BuddyUpdatedStatusEventBody struct {
	CharacterId   uint32 `json:"characterId"`
	Group         string `json:"group"`
	CharacterName string `json:"characterName"`
	ChannelId     int8   `json:"channelId"`
	InShop        bool   `json:"inShop"`
}

type BuddyChannelChangeStatusEventBody struct {
	CharacterId uint32 `json:"characterId"`
	ChannelId   int8   `json:"channelId"`
}

type BuddyCapacityChangeStatusEventBody struct {
	Capacity byte `json:"capacity"`
}

type ErrorStatusEventBody struct {
	Error string `json:"error"`
}
