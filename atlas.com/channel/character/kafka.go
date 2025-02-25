package character

const (
	EnvCommandTopic            = "COMMAND_TOPIC_CHARACTER"
	CommandRequestDistributeAp = "REQUEST_DISTRIBUTE_AP"
	CommandRequestDistributeSp = "REQUEST_DISTRIBUTE_SP"
	CommandRequestDropMeso     = "REQUEST_DROP_MESO"
	CommandChangeHP            = "CHANGE_HP"
	CommandChangeMP            = "CHANGE_MP"

	CommandDistributeApAbilityStrength     = "STRENGTH"
	CommandDistributeApAbilityDexterity    = "DEXTERITY"
	CommandDistributeApAbilityIntelligence = "INTELLIGENCE"
	CommandDistributeApAbilityLuck         = "LUCK"
	CommandDistributeApAbilityHp           = "HP"
	CommandDistributeApAbilityMp           = "MP"

	EnvCommandTopicMovement = "COMMAND_TOPIC_CHARACTER_MOVEMENT"

	MovementTypeNormal        = "NORMAL"
	MovementTypeTeleport      = "TELEPORT"
	MovementTypeStartFallDown = "START_FALL_DOWN"
	MovementTypeFlyingBlock   = "FLYING_BLOCK"
	MovementTypeJump          = "JUMP"
	MovementTypeStatChange    = "STAT_CHANGE"
)

type command[E any] struct {
	WorldId     byte   `json:"worldId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
	Body        E      `json:"body"`
}

type DistributePair struct {
	Ability string `json:"ability"`
	Amount  int8   `json:"amount"`
}

type requestDistributeApCommandBody struct {
	Distributions []DistributePair `json:"distributions"`
}

type requestDistributeSpCommandBody struct {
	SkillId uint32 `json:"skilId"`
	Amount  int8   `json:"amount"`
}

type requestDropMesoCommandBody struct {
	ChannelId byte   `json:"channelId"`
	MapId     uint32 `json:"mapId"`
	Amount    uint32 `json:"amount"`
}

type changeHPCommandBody struct {
	ChannelId byte  `json:"channelId"`
	Amount    int16 `json:"amount"`
}

type changeMPCommandBody struct {
	ChannelId byte  `json:"channelId"`
	Amount    int16 `json:"amount"`
}

type movementCommand struct {
	WorldId     byte     `json:"worldId"`
	ChannelId   byte     `json:"channelId"`
	MapId       uint32   `json:"mapId"`
	CharacterId uint32   `json:"characterId"`
	Movement    movement `json:"movement"`
}

type movement struct {
	StartX   int16     `json:"startX"`
	StartY   int16     `json:"startY"`
	Elements []element `json:"elements"`
}

type element struct {
	TypeStr     string `json:"typeStr"`
	TypeVal     byte   `json:"typeVal"`
	StartX      int16  `json:"startX"`
	StartY      int16  `json:"startY"`
	MoveAction  byte   `json:"moveAction"`
	Stat        byte   `json:"stat"`
	X           int16  `json:"x"`
	Y           int16  `json:"y"`
	VX          int16  `json:"vX"`
	VY          int16  `json:"vY"`
	FH          int16  `json:"fh"`
	FHFallStart int16  `json:"fhFallStart"`
	XOffset     int16  `json:"xOffset"`
	YOffset     int16  `json:"yOffset"`
	TimeElapsed int16  `json:"timeElapsed"`
}
