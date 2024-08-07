package writer

import (
	"atlas-channel/socket/model"
	"atlas-channel/tenant"
	"github.com/Chronicle20/atlas-socket/request"
	"github.com/Chronicle20/atlas-socket/response"
	"github.com/sirupsen/logrus"
	"time"
)

const MoveMonster = "MoveMonster"

func MoveMonsterBody(l logrus.FieldLogger, t tenant.Model) func(uniqueId uint32, bNotForceLandingWhenDiscard bool, bNotChangeAction bool, bNextAttackPossible bool, bLeft int8, skillId int16, skillLevel int16, multiTargets model.MultiTargetForBall, randTimeForAreaAttack model.RandTimeForAreaAttack, rawMovement []byte) BodyProducer {
	return func(uniqueId uint32, bNotForceLandingWhenDiscard bool, bNotChangeAction bool, bNextAttackPossible bool, bLeft int8, skillId int16, skillLevel int16, multiTargets model.MultiTargetForBall, randTimeForAreaAttack model.RandTimeForAreaAttack, rawMovement []byte) BodyProducer {
		return func(w *response.Writer, options map[string]interface{}) []byte {
			md := model.Movement{}
			r := request.NewRequestReader((*request.Request)(&rawMovement), time.Now().Unix())
			md.Decode(l, t, options)(&r)

			w.WriteInt(uniqueId)
			w.WriteBool(bNotForceLandingWhenDiscard)
			if (t.Region == "GMS" && t.MajorVersion > 83) || t.Region == "JMS" {
				w.WriteBool(bNotChangeAction)
			}
			w.WriteBool(bNextAttackPossible)
			w.WriteInt8(bLeft)
			w.WriteInt16(skillId)
			w.WriteInt16(skillLevel)
			if (t.Region == "GMS" && t.MajorVersion > 83) || t.Region == "JMS" {
				multiTargets.Encode(l, t, options)(w)
				randTimeForAreaAttack.Encode(l, t, options)(w)
			}
			md.Encode(l, t, options)(w)
			return w.Bytes()
		}
	}
}
