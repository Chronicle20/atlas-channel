package session

import (
	session2 "atlas-channel/kafka/message/session"
	"github.com/Chronicle20/atlas-constants/channel"
	"github.com/Chronicle20/atlas-constants/world"
	"github.com/Chronicle20/atlas-kafka/producer"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func StatusEventProvider(sessionId uuid.UUID, accountId uint32, characterId uint32, worldId world.Id, channelId channel.Id, eventType string) model.Provider[[]kafka.Message] {
	key := producer.CreateKey(int(characterId))
	value := &session2.StatusEvent{
		SessionId:   sessionId,
		AccountId:   accountId,
		CharacterId: characterId,
		WorldId:     byte(worldId),
		ChannelId:   byte(channelId),
		Issuer:      session2.EventSessionStatusIssuerChannel,
		Type:        eventType,
	}
	return producer.SingleMessageProvider(key, value)
}

func CreatedStatusEventProvider(sessionId uuid.UUID, accountId uint32, characterId uint32, worldId world.Id, channelId channel.Id) model.Provider[[]kafka.Message] {
	return StatusEventProvider(sessionId, accountId, characterId, worldId, channelId, session2.EventSessionStatusTypeCreated)
}

func DestroyedStatusEventProvider(sessionId uuid.UUID, accountId uint32, characterId uint32, worldId world.Id, channelId channel.Id) model.Provider[[]kafka.Message] {
	return StatusEventProvider(sessionId, accountId, characterId, worldId, channelId, session2.EventSessionStatusTypeDestroyed)
}
