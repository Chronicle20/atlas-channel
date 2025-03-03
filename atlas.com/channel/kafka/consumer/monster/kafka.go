package monster

const (
	EnvEventTopicStatus   = "EVENT_TOPIC_MONSTER_STATUS"
	EnvEventTopicMovement = "EVENT_TOPIC_MONSTER_MOVEMENT"

	EventStatusCreated      = "CREATED"
	EventStatusDestroyed    = "DESTROYED"
	EventStatusStartControl = "START_CONTROL"
	EventStatusStopControl  = "STOP_CONTROL"
	EventStatusDamaged      = "DAMAGED"
	EventStatusKilled       = "KILLED"

	MovementTypeNormal        = "NORMAL"
	MovementTypeTeleport      = "TELEPORT"
	MovementTypeStartFallDown = "START_FALL_DOWN"
	MovementTypeFlyingBlock   = "FLYING_BLOCK"
	MovementTypeJump          = "JUMP"
	MovementTypeStatChange    = "STAT_CHANGE"
)

type statusEvent[E any] struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	MapId     uint32 `json:"mapId"`
	UniqueId  uint32 `json:"uniqueId"`
	MonsterId uint32 `json:"monsterId"`
	Type      string `json:"type"`
	Body      E      `json:"body"`
}

type statusEventCreatedBody struct {
	ActorId uint32 `json:"actorId"`
}

type statusEventDestroyedBody struct {
	ActorId uint32 `json:"actorId"`
}

type statusEventStartControlBody struct {
	ActorId uint32 `json:"actorId"`
	X       int16  `json:"x"`
	Y       int16  `json:"y"`
	Stance  byte   `json:"stance"`
	FH      int16  `json:"fh"`
	Team    int8   `json:"team"`
}

type statusEventStopControlBody struct {
	ActorId uint32 `json:"actorId"`
}

type statusEventDamagedBody struct {
	X             int16         `json:"x"`
	Y             int16         `json:"y"`
	ActorId       uint32        `json:"actorId"`
	DamageEntries []damageEntry `json:"damageEntries"`
}

type statusEventKilledBody struct {
	X             int16         `json:"x"`
	Y             int16         `json:"y"`
	ActorId       uint32        `json:"actorId"`
	DamageEntries []damageEntry `json:"damageEntries"`
}

type damageEntry struct {
	CharacterId uint32 `json:"characterId"`
	Damage      int64  `json:"damage"`
}

type movementEvent struct {
	WorldId       byte       `json:"worldId"`
	ChannelId     byte       `json:"channelId"`
	MapId         uint32     `json:"mapId"`
	UniqueId      uint32     `json:"uniqueId"`
	ObserverId    uint32     `json:"observerId"`
	SkillPossible bool       `json:"skillPossible"`
	Skill         int8       `json:"skill"`
	SkillId       int16      `json:"skillId"`
	SkillLevel    int16      `json:"skillLevel"`
	MultiTarget   []position `json:"multiTarget"`
	RandomTimes   []int32    `json:"randomTimes"`
	Movement      movement   `json:"movement"`
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

type position struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}
