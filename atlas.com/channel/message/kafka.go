package message

const (
	EnvCommandTopicChat = "COMMAND_TOPIC_CHARACTER_CHAT"

	ChatTypeGeneral  = "GENERAL"
	ChatTypeBuddy    = "BUDDY"
	ChatTypeParty    = "PARTY"
	ChatTypeGuild    = "GUILD"
	ChatTypeAlliance = "Alliance"
)

type chatCommand[E any] struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	CharacterId uint32 `json:"characterId"`
	Message     string `json:"message"`
	Type        string `json:"type"`
	Body        E      `json:"body"`
}

type generalChatBody struct {
	BalloonOnly bool `json:"balloonOnly"`
}

type multiChatBody struct {
	Recipients []uint32 `json:"recipients"`
}
