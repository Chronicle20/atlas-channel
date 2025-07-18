package effect

import (
	"atlas-channel/data/skill/effect/statup"
	"github.com/Chronicle20/atlas-model/model"
)

type RestModel struct {
	WeaponAttack      int16   `json:"weaponAttack"`
	MagicAttack       int16   `json:"magicAttack"`
	WeaponDefense     int16   `json:"weaponDefense"`
	MagicDefense      int16   `json:"magicDefense"`
	Accuracy          int16   `json:"accuracy"`
	Avoidability      int16   `json:"avoidability"`
	Speed             int16   `json:"speed"`
	Jump              int16   `json:"jump"`
	HP                uint16  `json:"hp"`
	MP                uint16  `json:"mp"`
	HPR               float64 `json:"hpR"`
	MPR               float64 `json:"mpR"`
	MHPRRate          uint16  `json:"MHPRRate"`
	MMPRRate          uint16  `json:"MMPRRate"`
	MobSkill          uint16  `json:"mobSkill"`
	MobSkillLevel     uint16  `json:"mobSkillLevel"`
	MHPR              byte    `json:"mhpr"`
	MMPR              byte    `json:"mmpr"`
	HPConsume         uint16  `json:"HPConsume"`
	MPConsume         uint16  `json:"MPConsume"`
	Duration          int32   `json:"duration"`
	Target            uint32  `json:"target"`
	Barrier           int32   `json:"barrier"`
	Mob               uint32  `json:"mob"`
	OverTime          bool    `json:"overTime"`
	RepeatEffect      bool    `json:"repeatEffect"`
	MoveTo            int32   `json:"moveTo"`
	CP                uint32  `json:"cp"`
	NuffSkill         uint32  `json:"nuffSkill"`
	Skill             bool    `json:"skill"`
	X                 int16   `json:"x"`
	Y                 int16   `json:"y"`
	MobCount          uint32  `json:"mobCount"`
	MoneyConsume      uint32  `json:"moneyConsume"`
	Cooldown          uint32  `json:"cooldown"`
	MorphId           uint32  `json:"morphId"`
	Ghost             uint32  `json:"ghost"`
	Fatigue           uint32  `json:"fatigue"`
	Berserk           uint32  `json:"berserk"`
	Booster           uint32  `json:"booster"`
	Prop              float64 `json:"prop"`
	ItemConsume       uint32  `json:"itemConsume"`
	ItemConsumeAmount uint32  `json:"itemConsumeAmount"`
	Damage            uint32  `json:"damage"`
	AttackCount       uint32  `json:"attackCount"`
	FixDamage         int32   `json:"fixDamage"`
	//LT Point
	//RB Point
	BulletCount          uint16             `json:"bulletCount"`
	BulletConsume        uint16             `json:"bulletConsume"`
	MapProtection        byte               `json:"mapProtection"`
	CureAbnormalStatuses []string           `json:"cureAbnormalStatuses"`
	Statups              []statup.RestModel `json:"statups"`
	MonsterStatus        map[string]uint32  `json:"monsterStatus"`
	CardStats            cardItemUp         `json:"cardStats"`
}

func Extract(rm RestModel) (Model, error) {
	su, err := model.SliceMap(statup.Extract)(model.FixedProvider(rm.Statups))()()
	if err != nil {
		return Model{}, err
	}

	return Model{
		weaponAttack:         rm.WeaponAttack,
		magicAttack:          rm.MagicAttack,
		weaponDefense:        rm.WeaponDefense,
		magicDefense:         rm.MagicDefense,
		accuracy:             rm.Accuracy,
		avoidability:         rm.Avoidability,
		speed:                rm.Speed,
		jump:                 rm.Jump,
		hp:                   rm.HP,
		mp:                   rm.MP,
		hpr:                  rm.HPR,
		mpr:                  rm.MPR,
		mhprRate:             rm.MHPRRate,
		mmprRate:             rm.MMPRRate,
		mobSkill:             rm.MobSkill,
		mobSkillLevel:        rm.MobSkillLevel,
		mhpR:                 rm.MHPR,
		mmpR:                 rm.MMPR,
		hpCon:                rm.HPConsume,
		mpCon:                rm.MPConsume,
		duration:             rm.Duration,
		target:               rm.Target,
		barrier:              rm.Barrier,
		mob:                  rm.Mob,
		overtime:             rm.OverTime,
		repeatEffect:         rm.RepeatEffect,
		moveTo:               rm.MoveTo,
		cp:                   rm.CP,
		nuffSkill:            rm.NuffSkill,
		skill:                rm.Skill,
		x:                    rm.X,
		y:                    rm.Y,
		mobCount:             rm.MobCount,
		moneyCon:             rm.MoneyConsume,
		cooldown:             rm.Cooldown,
		morphId:              rm.MorphId,
		ghost:                rm.Ghost,
		fatigue:              rm.Fatigue,
		berserk:              rm.Berserk,
		booster:              rm.Booster,
		prop:                 rm.Prop,
		itemCon:              rm.ItemConsume,
		itemConNo:            rm.ItemConsumeAmount,
		damage:               rm.Damage,
		attackCount:          rm.AttackCount,
		fixDamage:            rm.FixDamage,
		bulletCount:          rm.BulletCount,
		bulletConsume:        rm.BulletConsume,
		mapProtection:        rm.MapProtection,
		cureAbnormalStatuses: rm.CureAbnormalStatuses,
		statups:              su,
		monsterStatus:        rm.MonsterStatus,
	}, nil
}

type monsterStatus struct {
	Status string `json:"status"`
	Value  uint32 `json:"value"`
}

type cardItemUp struct {
	ItemCode    uint32 `json:"itemCode"`
	Probability uint32 `json:"probability"`
	Areas       []area `json:"areas"`
	InParty     bool   `json:"inParty"`
}

type area struct {
	Start uint32 `json:"start"`
	End   uint32 `json:"end"`
}
