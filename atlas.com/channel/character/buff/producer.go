package buff

import (
	"atlas-channel/skill/effect/statup"
	"github.com/Chronicle20/atlas-kafka/producer"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/segmentio/kafka-go"
)

func applyCommandProvider(worldId byte, channelId byte, characterId uint32, fromId uint32, sourceId uint32, duration int32, statups []statup.Model) model.Provider[[]kafka.Message] {
	changes := make([]statChange, 0)
	for _, su := range statups {
		changes = append(changes, statChange{
			Type:   su.Mask(),
			Amount: su.Amount(),
		})
	}

	key := producer.CreateKey(int(characterId))
	value := &command[applyCommandBody]{
		WorldId:     worldId,
		ChannelId:   channelId,
		CharacterId: characterId,
		Type:        CommandTypeApply,
		Body: applyCommandBody{
			FromId:   fromId,
			SourceId: sourceId,
			Duration: duration,
			Changes:  changes,
		},
	}
	return producer.SingleMessageProvider(key, value)
}

func cancelCommandProvider(worldId byte, channelId byte, characterId uint32, sourceId uint32) model.Provider[[]kafka.Message] {
	key := producer.CreateKey(int(characterId))
	value := &command[cancelCommandBody]{
		WorldId:     worldId,
		ChannelId:   channelId,
		CharacterId: characterId,
		Type:        CommandTypeCancel,
		Body: cancelCommandBody{
			SourceId: sourceId,
		},
	}
	return producer.SingleMessageProvider(key, value)
}
