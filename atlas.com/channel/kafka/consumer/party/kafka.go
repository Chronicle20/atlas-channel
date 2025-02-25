package party

const (
	EnvEventStatusTopic                             = "EVENT_TOPIC_PARTY_STATUS"
	EventPartyStatusTypeCreated                     = "CREATED"
	EventPartyStatusTypeJoined                      = "JOINED"
	EventPartyStatusTypeLeft                        = "LEFT"
	EventPartyStatusTypeExpel                       = "EXPEL"
	EventPartyStatusTypeDisband                     = "DISBAND"
	EventPartyStatusTypeChangeLeader                = "CHANGE_LEADER"
	EventPartyStatusTypeError                       = "ERROR"
)

type statusEvent[E any] struct {
	ActorId uint32 `json:"actorId"`
	WorldId byte   `json:"worldId"`
	PartyId uint32 `json:"partyId"`
	Type    string `json:"type"`
	Body    E      `json:"body"`
}

type createdEventBody struct {
}

type joinedEventBody struct {
}

type leftEventBody struct {
}

type expelEventBody struct {
	CharacterId uint32 `json:"characterId"`
}

type disbandEventBody struct {
	Members []uint32 `json:"members"`
}

type changeLeaderEventBody struct {
	CharacterId  uint32 `json:"characterId"`
	Disconnected bool   `json:"disconnected"`
}

type errorEventBody struct {
	Type          string `json:"type"`
	CharacterName string `json:"characterName"`
}
