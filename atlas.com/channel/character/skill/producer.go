package skill

import (
	"github.com/Chronicle20/atlas-kafka/producer"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/segmentio/kafka-go"
)

func setCooldownCommandProvider(characterId uint32, id uint32, cooldown uint32) model.Provider[[]kafka.Message] {
	key := producer.CreateKey(int(characterId))
	value := &command[setCooldownBody]{
		CharacterId: characterId,
		Type:        CommandTypeSetCooldown,
		Body: setCooldownBody{
			SkillId:  id,
			Cooldown: cooldown,
		},
	}
	return producer.SingleMessageProvider(key, value)
}
