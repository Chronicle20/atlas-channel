package handler

import (
	"atlas-channel/account"
	"atlas-channel/buddylist"
	"atlas-channel/cashshop"
	"atlas-channel/character"
	"atlas-channel/session"
	"atlas-channel/socket/writer"
	"context"
	"github.com/Chronicle20/atlas-socket/request"
	tenant "github.com/Chronicle20/atlas-tenant"
	"github.com/sirupsen/logrus"
)

const CashShopEntryHandle = "CashShopEntryHandle"

func CashShopEntryHandleFunc(l logrus.FieldLogger, ctx context.Context, wp writer.Producer) func(s session.Model, r *request.Reader, readerOptions map[string]interface{}) {
	t := tenant.MustFromContext(ctx)
	return func(s session.Model, r *request.Reader, readerOptions map[string]interface{}) {
		updateTime := r.ReadUint32()
		l.Debugf("Character [%d] is attempting to enter the cash shop. update_time [%d].", s.CharacterId(), updateTime)

		// TODO block when performing vega scrolling
		// TODO block when in event
		// TODO block when in mini dungeon
		// TODO block when already in cash shop

		a, err := account.GetById(l)(ctx)(s.AccountId())
		if err != nil {
			l.WithError(err).Errorf("Unable to locate account [%d] attempting to enter cash shop.", s.AccountId())
			_ = session.Destroy(l, ctx, session.GetRegistry())(s)
			return
		}
		c, err := character.GetByIdWithInventory(l)(ctx)(character.SkillModelDecorator(l)(ctx))(s.CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Unable to locate character [%d] attempting to enter cash shop.", s.CharacterId())
			_ = session.Destroy(l, ctx, session.GetRegistry())(s)
			return
		}
		bl, err := buddylist.GetById(l)(ctx)(s.CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Unable to locate buddylist [%d] attempting to enter cash shop.", s.CharacterId())
			_ = session.Destroy(l, ctx, session.GetRegistry())(s)
			return
		}

		err = session.Announce(l)(ctx)(wp)(writer.CashShopOpen)(writer.CashShopOpenBody(l)(t, a, c, bl))(s)
		if err != nil {
			return
		}

		//err = session.Announce(l)(wp)(writer.CashShopOperation)(s, writer.CashShopCashInventoryBody(l)(s.Tenant()))
		//if err != nil {
		//	return
		//}
		//
		//err = session.Announce(l)(wp)(writer.CashShopOperation)(s, writer.CashShopCashGiftsBody(l)(s.Tenant()))
		//if err != nil {
		//	return
		//}
		//
		//err = session.Announce(l)(wp)(writer.CashShopOperation)(s, writer.CashShopWishListBody(l)(s.Tenant()))
		//if err != nil {
		//	return
		//}

		err = session.Announce(l)(ctx)(wp)(writer.CashShopCashQueryResult)(writer.CashShopCashQueryResultBody(l)(t))(s)
		if err != nil {
			return
		}

		err = cashshop.Enter(l)(ctx)(s.CharacterId(), s.Map())
		if err != nil {
			l.WithError(err).Errorf("Unable to announce [%d] has entered the cash shop.", s.CharacterId())
		}
	}
}
