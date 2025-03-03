package writer

import (
	"atlas-channel/buddylist"
	"atlas-channel/character"
	"atlas-channel/character/equipment/slot"
	"atlas-channel/character/inventory/equipable"
	"atlas-channel/character/inventory/item"
	"github.com/Chronicle20/atlas-constants/channel"
	_map "github.com/Chronicle20/atlas-constants/map"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/Chronicle20/atlas-socket/response"
	"github.com/Chronicle20/atlas-tenant"
	"github.com/sirupsen/logrus"
	"math"
	"math/rand"
	"time"
)

const SetField = "SetField"

func WarpToMapBody(l logrus.FieldLogger, tenant tenant.Model) func(channelId channel.Id, mapId _map.Id, portalId uint32, hp uint16) BodyProducer {
	return func(channelId channel.Id, mapId _map.Id, portalId uint32, hp uint16) BodyProducer {
		return func(w *response.Writer, options map[string]interface{}) []byte {
			if (tenant.Region() == "GMS" && tenant.MajorVersion() > 83) || tenant.Region() == "JMS" {
				w.WriteShort(0) // decode opt, loop with 2 decode 4s
			}
			w.WriteInt(uint32(channelId))
			if tenant.Region() == "JMS" {
				w.WriteByte(0)
				w.WriteInt(0)
			}
			w.WriteByte(0) // sNotifierMessage
			w.WriteByte(0) // bCharacterData
			if (tenant.Region() == "GMS" && tenant.MajorVersion() > 28) || tenant.Region() == "JMS" {
				w.WriteShort(0) // nNotifierCheck
				w.WriteByte(0)  // revive
			}
			w.WriteInt(uint32(mapId))
			w.WriteByte(byte(portalId))
			w.WriteShort(hp)
			if tenant.Region() == "GMS" && tenant.MajorVersion() > 28 {
				w.WriteBool(false) // Chasing?
				if false {
					w.WriteInt(0)
					w.WriteInt(0)
				}
			}
			w.WriteLong(uint64(getTime(timeNow())))
			return w.Bytes()
		}
	}
}

func SetFieldBody(l logrus.FieldLogger, tenant tenant.Model) func(channelId channel.Id, c character.Model, bl buddylist.Model) BodyProducer {
	return func(channelId channel.Id, c character.Model, bl buddylist.Model) BodyProducer {
		return func(w *response.Writer, options map[string]interface{}) []byte {
			if (tenant.Region() == "GMS" && tenant.MajorVersion() > 83) || tenant.Region() == "JMS" {
				w.WriteShort(0) // decode opt, loop with 2 decode 4s
			}
			w.WriteInt(uint32(channelId))
			if tenant.Region() == "JMS" {
				w.WriteByte(0)
				w.WriteInt(0)
			}
			w.WriteByte(1) // sNotifierMessage
			w.WriteByte(1) // bCharacterData

			var seedSize = 3
			if (tenant.Region() == "GMS" && tenant.MajorVersion() > 28) || tenant.Region() == "JMS" {
				w.WriteShort(0) // nNotifierCheck, if non zero STRs are encoded
			} else {
				seedSize = 4
			}

			// damage seed
			for i := 0; i < seedSize; i++ {
				w.WriteInt(rand.Uint32())
			}

			WriteCharacterInfo(tenant)(w)(c, bl)
			if (tenant.Region() == "GMS" && tenant.MajorVersion() > 83) || tenant.Region() == "JMS" {
				w.WriteInt(0) // logout gifts
				w.WriteInt(0)
				w.WriteInt(0)
				w.WriteInt(0)
			}
			w.WriteInt64(getTime(timeNow()))
			return w.Bytes()
		}
	}
}

func WriteCharacterInfo(tenant tenant.Model) func(w *response.Writer) func(c character.Model, bl buddylist.Model) {
	return func(w *response.Writer) func(c character.Model, bl buddylist.Model) {
		return func(c character.Model, bl buddylist.Model) {
			if (tenant.Region() == "GMS" && tenant.MajorVersion() > 28) || tenant.Region() == "JMS" {
				w.WriteInt64(-1) // dbcharFlag
				w.WriteByte(0)   // something about SN, I believe this is size of list
			} else {
				w.WriteInt16(-1) // dbcharFlag
			}

			WriteCharacterStatistics(tenant)(w, c)
			w.WriteByte(bl.Capacity())

			if (tenant.Region() == "GMS" && tenant.MajorVersion() > 28) || tenant.Region() == "JMS" {
				if true {
					w.WriteByte(0)
				} else {
					w.WriteByte(1)
					w.WriteAsciiString("") // linked name
				}
			}
			w.WriteInt(c.Meso())

			if tenant.Region() == "JMS" {
				w.WriteInt(c.Id())
				w.WriteInt(0) // dama / gachapon items
				w.WriteInt(0)
			}
			WriteInventoryInfo(tenant)(w, c)
			WriteSkillInfo(tenant)(w, c)
			WriteQuestInfo(tenant)(w, c)
			WriteMiniGameInfo(tenant)(w, c)
			WriteRingInfo(tenant)(w, c)
			WriteTeleportInfo(tenant)(w, c)
			if tenant.Region() == "JMS" {
				w.WriteShort(0)
			}

			if (tenant.Region() == "GMS" && tenant.MajorVersion() > 28) || tenant.Region() == "JMS" {
				WriteMonsterBookInfo(tenant)(w, c)
				if tenant.Region() == "GMS" {
					WriteNewYearInfo(tenant)(w, c)
					WriteAreaInfo(tenant)(w, c)
				} else if tenant.Region() == "JMS" {
					w.WriteShort(0)
				}
				w.WriteShort(0)
			}
		}
	}
}

func WriteAreaInfo(tenant tenant.Model) func(w *response.Writer, c character.Model) {
	return func(w *response.Writer, c character.Model) {
		w.WriteShort(0)
	}
}

func WriteNewYearInfo(tenant tenant.Model) func(w *response.Writer, c character.Model) {
	return func(w *response.Writer, c character.Model) {
		w.WriteShort(0)
	}
}

func WriteMonsterBookInfo(tenant tenant.Model) func(w *response.Writer, c character.Model) {
	return func(w *response.Writer, c character.Model) {
		w.WriteInt(0) // cover id
		w.WriteByte(0)
		w.WriteShort(0) // card size
	}
}

func WriteTeleportInfo(tenant tenant.Model) func(w *response.Writer, c character.Model) {
	return func(w *response.Writer, c character.Model) {
		for i := 0; i < 5; i++ {
			w.WriteInt(uint32(_map.EmptyMapId))
		}

		if (tenant.Region() == "GMS" && tenant.MajorVersion() > 28) || tenant.Region() == "JMS" {
			for j := 0; j < 10; j++ {
				w.WriteInt(uint32(_map.EmptyMapId))
			}
		}
	}
}

func WriteRingInfo(tenant tenant.Model) func(w *response.Writer, c character.Model) {
	return func(w *response.Writer, c character.Model) {
		w.WriteShort(0) // crush rings

		if (tenant.Region() == "GMS" && tenant.MajorVersion() > 28) || tenant.Region() == "JMS" {
			w.WriteShort(0) // friendship rings
			w.WriteShort(0) // partner
		}
	}
}

func WriteMiniGameInfo(tenant tenant.Model) func(w *response.Writer, c character.Model) {
	return func(w *response.Writer, c character.Model) {
		w.WriteShort(0)
	}
}

func WriteQuestInfo(tenant tenant.Model) func(w *response.Writer, c character.Model) {
	return func(w *response.Writer, c character.Model) {
		w.WriteShort(0) // started size
		if tenant.Region() == "JMS" {
			w.WriteShort(0)
		}

		if (tenant.Region() == "GMS" && tenant.MajorVersion() > 12) || tenant.Region() == "JMS" {
			w.WriteShort(0) // completed size
		}
	}
}

func WriteSkillInfo(tenant tenant.Model) func(w *response.Writer, c character.Model) {
	return func(w *response.Writer, c character.Model) {
		var onCooldown []int

		w.WriteShort(uint16(len(c.Skills())))
		for i, s := range c.Skills() {
			w.WriteInt(s.Id())
			w.WriteInt(uint32(s.Level()))
			w.WriteInt64(msTime(s.Expiration()))
			if s.IsFourthJob() {
				w.WriteInt(uint32(s.MasterLevel()))
			}
			if s.OnCooldown() {
				onCooldown = append(onCooldown, i)
			}
		}

		if (tenant.Region() == "GMS" && tenant.MajorVersion() > 28) || tenant.Region() == "JMS" {
			w.WriteShort(uint16(len(onCooldown)))
			for _, i := range onCooldown {
				s := c.Skills()[i]
				w.WriteInt(s.Id())
				cd := uint32(s.CooldownExpiresAt().Sub(time.Now()).Seconds())
				w.WriteShort(uint16(cd))
			}
		}
	}
}

const (
	DefaultTime int64 = 150842304000000000
	ZeroTime    int64 = 94354848000000000
	Permanent   int64 = 150841440000000000
)

func getTime(utcTimestamp int64) int64 {
	if utcTimestamp < 0 && utcTimestamp >= -3 {
		if utcTimestamp == -1 {
			return DefaultTime //high number ll
		} else if utcTimestamp == -2 {
			return ZeroTime
		} else {
			return Permanent
		}
	}

	ftUtOffset := 116444736010800000 + (10000 * timeNow())
	return utcTimestamp*10000 + ftUtOffset
}

func timeNow() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func WriteInventoryInfo(tenant tenant.Model) func(w *response.Writer, character character.Model) {
	return func(w *response.Writer, character character.Model) {

		if (tenant.Region() == "GMS" && tenant.MajorVersion() > 12) || tenant.Region() == "JMS" {
			w.WriteByte(byte(character.Inventory().Equipable().Capacity()))
			w.WriteByte(byte(character.Inventory().Use().Capacity()))
			w.WriteByte(byte(character.Inventory().Setup().Capacity()))
			w.WriteByte(byte(character.Inventory().Etc().Capacity()))
			w.WriteByte(byte(character.Inventory().Cash().Capacity()))
		}

		if (tenant.Region() == "GMS" && tenant.MajorVersion() > 28) || tenant.Region() == "JMS" {
			w.WriteLong(uint64(getTime(-2)))
		}

		// regular equipment
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Hat())
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Medal())
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Forehead())
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Ring1())
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Ring2())
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Eye())
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Earring())
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Shoulder())
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Cape())
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Top())
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Pendant())
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Weapon())
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Shield())
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Gloves())
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Bottom())
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Belt())
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Ring3())
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Ring4())
		WriteEquipableIfPresent(tenant)(w, character.Equipment().Shoes())
		if (tenant.Region() == "GMS" && tenant.MajorVersion() > 28) || tenant.Region() == "JMS" {
			w.WriteShort(0)
		} else {
			w.WriteByte(0)
		}
		// cash equipment
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Hat())
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Medal())
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Forehead())
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Ring1())
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Ring2())
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Eye())
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Earring())
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Shoulder())
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Cape())
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Top())
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Pendant())
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Weapon())
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Shield())
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Gloves())
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Bottom())
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Belt())
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Ring3())
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Ring4())
		WriteCashEquipableIfPresent(tenant)(w, character.Equipment().Shoes())
		if (tenant.Region() == "GMS" && tenant.MajorVersion() > 28) || tenant.Region() == "JMS" {
			w.WriteShort(0)
		} else {
			w.WriteByte(0)
		}

		// equipable inventory
		if tenant.Region() == "GMS" && tenant.MajorVersion() < 28 {
			w.WriteByte(byte(character.Inventory().Equipable().Capacity()))
		}
		_ = model.ForEachSlice(model.FixedProvider(character.Inventory().Equipable().Items()), WriteEquipableInfo(tenant)(w, false))
		if (tenant.Region() == "GMS" && tenant.MajorVersion() > 28) || tenant.Region() == "JMS" {
			w.WriteInt(0)
		} else {
			w.WriteByte(0)
		}

		// use inventory
		if tenant.Region() == "GMS" && tenant.MajorVersion() < 28 {
			w.WriteByte(byte(character.Inventory().Use().Capacity()))
		}
		_ = model.ForEachSlice(model.FixedProvider(character.Inventory().Use().Items()), WriteItemInfo(tenant)(w, false))
		w.WriteByte(0)

		// setup inventory
		if tenant.Region() == "GMS" && tenant.MajorVersion() < 28 {
			w.WriteByte(byte(character.Inventory().Setup().Capacity()))
		}
		_ = model.ForEachSlice(model.FixedProvider(character.Inventory().Setup().Items()), WriteItemInfo(tenant)(w, false))
		w.WriteByte(0)

		// etc inventory
		if tenant.Region() == "GMS" && tenant.MajorVersion() < 28 {
			w.WriteByte(byte(character.Inventory().Etc().Capacity()))
		}
		_ = model.ForEachSlice(model.FixedProvider(character.Inventory().Etc().Items()), WriteItemInfo(tenant)(w, false))
		w.WriteByte(0)

		// cash inventory
		if tenant.Region() == "GMS" && tenant.MajorVersion() < 28 {
			w.WriteByte(byte(character.Inventory().Cash().Capacity()))
		}
		_ = model.ForEachSlice(model.FixedProvider(character.Inventory().Cash().Items()), WriteCashItemInfo(tenant)(w, false))
		w.WriteByte(0)
	}
}

func WriteCashEquipableIfPresent(tenant tenant.Model) func(w *response.Writer, model slot.Model) {
	return func(w *response.Writer, model slot.Model) {
		if model.CashEquipable != nil {
			_ = WriteCashEquipableInfo(tenant)(w, false)(*model.CashEquipable)
		}
	}
}

func WriteCashEquipableInfo(tenant tenant.Model) func(w *response.Writer, zeroPosition bool) model.Operator[equipable.Model] {
	return func(w *response.Writer, zeroPosition bool) model.Operator[equipable.Model] {
		return func(e equipable.Model) error {
			slot := e.Slot()
			if !zeroPosition {
				slot = int16(math.Abs(float64(slot)))
				if slot > 100 {
					slot -= 100
				}
				if (tenant.Region() == "GMS" && tenant.MajorVersion() > 28) || tenant.Region() == "JMS" {
					w.WriteShort(uint16(slot))
				} else {
					w.WriteByte(byte(slot))
				}
			}

			if (tenant.Region() == "GMS" && tenant.MajorVersion() > 12) || tenant.Region() == "JMS" {
				w.WriteByte(1)
			}
			w.WriteInt(e.ItemId())
			w.WriteBool(true)
			if true {
				w.WriteLong(uint64(e.Id())) // cash sn
			}
			w.WriteLong(uint64(getTime(e.Expiration())))
			w.WriteByte(byte(e.Slots()))
			w.WriteByte(e.Level())
			if tenant.Region() == "JMS" {
				w.WriteByte(0)
			}
			w.WriteShort(e.Strength())
			w.WriteShort(e.Dexterity())
			w.WriteShort(e.Intelligence())
			w.WriteShort(e.Luck())
			w.WriteShort(e.HP())
			w.WriteShort(e.MP())
			w.WriteShort(e.WeaponAttack())
			w.WriteShort(e.MagicAttack())
			w.WriteShort(e.WeaponDefense())
			w.WriteShort(e.MagicDefense())
			w.WriteShort(e.Accuracy())
			w.WriteShort(e.Avoidability())
			w.WriteShort(e.Hands())
			w.WriteShort(e.Speed())
			w.WriteShort(e.Jump())

			if (tenant.Region() == "GMS" && tenant.MajorVersion() > 12) || tenant.Region() == "JMS" {
				w.WriteAsciiString(e.OwnerName())
				w.WriteShort(e.Flags())

				if (tenant.Region() == "GMS" && tenant.MajorVersion() > 28) || tenant.Region() == "JMS" {
					for i := 0; i < 10; i++ {
						w.WriteByte(0x40)
					}
					w.WriteLong(uint64(getTime(-2)))
					w.WriteInt32(-1)
				}
			}
			return nil
		}
	}
}

func WriteCashItemInfo(_ tenant.Model) func(w *response.Writer, zeroPosition bool) model.Operator[item.Model] {
	return func(w *response.Writer, zeroPosition bool) model.Operator[item.Model] {
		return func(i item.Model) error {
			if !zeroPosition {
				w.WriteInt8(int8(i.Slot()))
			}
			w.WriteByte(2)
			w.WriteInt(i.ItemId())
			w.WriteBool(true)
			w.WriteLong(777000000) // pet id, ring id, or incremental number
			w.WriteLong(uint64(getTime(i.Expiration())))
			w.WriteShort(uint16(i.Quantity()))
			w.WriteAsciiString(i.Owner())
			w.WriteShort(i.Flag())
			return nil
		}
	}
}

func WriteItemInfo(_ tenant.Model) func(w *response.Writer, zeroPosition bool) model.Operator[item.Model] {
	return func(w *response.Writer, zeroPosition bool) model.Operator[item.Model] {
		return func(i item.Model) error {
			if !zeroPosition {
				w.WriteInt8(int8(i.Slot()))
			}
			w.WriteByte(2)
			w.WriteInt(i.ItemId())
			w.WriteBool(false)
			w.WriteLong(uint64(getTime(i.Expiration())))
			w.WriteShort(uint16(i.Quantity()))
			w.WriteAsciiString(i.Owner())
			w.WriteShort(i.Flag())
			if i.Rechargeable() {
				w.WriteLong(0)
			}
			return nil
		}
	}
}

func WriteEquipableIfPresent(tenant tenant.Model) func(w *response.Writer, model slot.Model) {
	return func(w *response.Writer, model slot.Model) {
		if model.Equipable != nil {
			_ = WriteEquipableInfo(tenant)(w, false)(*model.Equipable)
		}
	}
}

func WriteEquipableInfo(tenant tenant.Model) func(w *response.Writer, zeroPosition bool) model.Operator[equipable.Model] {
	return func(w *response.Writer, zeroPosition bool) model.Operator[equipable.Model] {
		return func(e equipable.Model) error {
			slot := e.Slot()
			if !zeroPosition {
				slot = int16(math.Abs(float64(slot)))
				if slot > 100 {
					slot -= 100
				}
				if (tenant.Region() == "GMS" && tenant.MajorVersion() > 28) || tenant.Region() == "JMS" {
					w.WriteShort(uint16(slot))
				} else {
					w.WriteByte(byte(slot))
				}
			}

			if (tenant.Region() == "GMS" && tenant.MajorVersion() > 12) || tenant.Region() == "JMS" {
				w.WriteByte(1)
			}
			w.WriteInt(e.ItemId())
			w.WriteBool(false)
			if false {
				w.WriteLong(0) // cash sn
			}
			w.WriteLong(uint64(getTime(e.Expiration())))
			w.WriteByte(byte(e.Slots()))
			w.WriteByte(e.Level())
			if tenant.Region() == "JMS" {
				w.WriteByte(0)
			}
			w.WriteShort(e.Strength())
			w.WriteShort(e.Dexterity())
			w.WriteShort(e.Intelligence())
			w.WriteShort(e.Luck())
			w.WriteShort(e.HP())
			w.WriteShort(e.MP())
			w.WriteShort(e.WeaponAttack())
			w.WriteShort(e.MagicAttack())
			w.WriteShort(e.WeaponDefense())
			w.WriteShort(e.MagicDefense())
			w.WriteShort(e.Accuracy())
			w.WriteShort(e.Avoidability())
			w.WriteShort(e.Hands())
			w.WriteShort(e.Speed())
			w.WriteShort(e.Jump())

			if (tenant.Region() == "GMS" && tenant.MajorVersion() > 12) || tenant.Region() == "JMS" {
				w.WriteAsciiString(e.OwnerName())
				w.WriteShort(e.Flags())

				if (tenant.Region() == "GMS" && tenant.MajorVersion() > 28) || tenant.Region() == "JMS" {
					w.WriteByte(0)
					w.WriteByte(0)   // item level
					w.WriteInt(0)    // expNibble
					w.WriteInt32(-1) // durability

					if tenant.Region() == "JMS" {
						w.WriteByte(0)
						w.WriteShort(0)
						w.WriteShort(0)
						w.WriteShort(0)
						w.WriteShort(0)
						w.WriteShort(0)
						w.WriteInt(0)
					}

					w.WriteLong(0)
					w.WriteLong(uint64(getTime(-2)))
					w.WriteInt32(-1)
				}
			}
			return nil
		}
	}
}

func WriteCharacterStatistics(tenant tenant.Model) func(w *response.Writer, character character.Model) {
	return func(w *response.Writer, character character.Model) {
		w.WriteInt(character.Id())

		name := character.Name()
		if len(name) > 13 {
			name = name[:13]
		}
		padSize := 13 - len(name)
		w.WriteByteArray([]byte(name))
		for i := 0; i < padSize; i++ {
			w.WriteByte(0x0)
		}

		w.WriteByte(character.Gender())
		w.WriteByte(character.SkinColor())
		w.WriteInt(character.Face())
		w.WriteInt(character.Hair())

		if (tenant.Region() == "GMS" && tenant.MajorVersion() > 28) || tenant.Region() == "JMS" {
			writeForEachPet(w, character.Pets(), writePetId, writeEmptyPetId)
		} else {
			if len(character.Pets()) > 0 {
				w.WriteLong(character.Pets()[0].Id()) // pet cash id
			} else {
				w.WriteLong(0)
			}
		}
		w.WriteByte(character.Level())
		w.WriteShort(character.JobId())
		w.WriteShort(character.Strength())
		w.WriteShort(character.Dexterity())
		w.WriteShort(character.Intelligence())
		w.WriteShort(character.Luck())
		w.WriteShort(character.Hp())
		w.WriteShort(character.MaxHp())
		w.WriteShort(character.Mp())
		w.WriteShort(character.MaxMp())
		w.WriteShort(character.Ap())

		if character.HasSPTable() {
			WriteRemainingSkillInfo(w, character)
		} else {
			w.WriteShort(character.RemainingSp())
		}

		w.WriteInt(character.Experience())
		w.WriteInt16(character.Fame())
		if (tenant.Region() == "GMS" && tenant.MajorVersion() > 28) || tenant.Region() == "JMS" {
			w.WriteInt(character.GachaponExperience())
		}
		w.WriteInt(character.MapId())
		w.WriteByte(character.SpawnPoint())

		if tenant.Region() == "GMS" {
			if tenant.MajorVersion() > 12 {
				w.WriteInt(0)
			} else {
				w.WriteInt64(0)
				w.WriteInt(0)
				w.WriteInt(0)
			}
			if tenant.MajorVersion() >= 87 {
				w.WriteShort(0) // nSubJob
			}
		} else if tenant.Region() == "JMS" {
			w.WriteShort(0)
			w.WriteLong(0)
			w.WriteInt(0)
			w.WriteInt(0)
			w.WriteInt(0)
		}
	}
}

func WriteRemainingSkillInfo(w *response.Writer, character character.Model) {

}
