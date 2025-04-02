package handler

import (
	"atlas-channel/cashshop/wallet"
	"atlas-channel/session"
	"atlas-channel/socket/writer"
	"context"
	"github.com/Chronicle20/atlas-socket/request"
	"github.com/Chronicle20/atlas-tenant"
	"github.com/sirupsen/logrus"
)

const CashShopCheckWalletHandle = "CashShopCheckWalletHandle"

func CashShopCheckWalletHandleFunc(l logrus.FieldLogger, ctx context.Context, wp writer.Producer) func(s session.Model, r *request.Reader, readerOptions map[string]interface{}) {
	t := tenant.MustFromContext(ctx)
	return func(s session.Model, r *request.Reader, readerOptions map[string]interface{}) {
		w, err := wallet.GetByCharacterId(l)(ctx)(s.CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve cash shop wallet for character [%d].", s.CharacterId())
			err = session.Announce(l)(ctx)(wp)(writer.CashShopCashQueryResult)(writer.CashShopCashQueryResultBody(t)(0, 0, 0))(s)
			if err != nil {
				l.WithError(err).Errorf("Unable to announce default cash shop wallet to character [%d].", s.CharacterId())
				return
			}
		} else {
			err = session.Announce(l)(ctx)(wp)(writer.CashShopCashQueryResult)(writer.CashShopCashQueryResultBody(t)(w.Credit(), w.Points(), w.Prepaid()))(s)
			if err != nil {
				l.WithError(err).Errorf("Unable to announce cash shop wallet to character [%d].", s.CharacterId())
				return
			}
		}
	}
}
