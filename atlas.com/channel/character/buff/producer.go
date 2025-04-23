package buff

import (
	"atlas-channel/kafka/message/buff"
	"atlas-channel/skill/effect/statup"
	_map "github.com/Chronicle20/atlas-constants/map"
	"github.com/Chronicle20/atlas-kafka/producer"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/segmentio/kafka-go"
)

func ApplyCommandProvider(m _map.Model, characterId uint32, fromId uint32, sourceId int32, duration int32, statups []statup.Model) model.Provider[[]kafka.Message] {
	changes := make([]buff.StatChange, 0)
	for _, su := range statups {
		changes = append(changes, buff.StatChange{
			Type:   su.Mask(),
			Amount: su.Amount(),
		})
	}

	key := producer.CreateKey(int(characterId))
	value := &buff.Command[buff.ApplyCommandBody]{
		WorldId:     byte(m.WorldId()),
		ChannelId:   byte(m.ChannelId()),
		CharacterId: characterId,
		Type:        buff.CommandTypeApply,
		Body: buff.ApplyCommandBody{
			FromId:   fromId,
			SourceId: sourceId,
			Duration: duration,
			Changes:  changes,
		},
	}
	return producer.SingleMessageProvider(key, value)
}

func CancelCommandProvider(m _map.Model, characterId uint32, sourceId int32) model.Provider[[]kafka.Message] {
	key := producer.CreateKey(int(characterId))
	value := &buff.Command[buff.CancelCommandBody]{
		WorldId:     byte(m.WorldId()),
		ChannelId:   byte(m.ChannelId()),
		CharacterId: characterId,
		Type:        buff.CommandTypeCancel,
		Body: buff.CancelCommandBody{
			SourceId: sourceId,
		},
	}
	return producer.SingleMessageProvider(key, value)
}
