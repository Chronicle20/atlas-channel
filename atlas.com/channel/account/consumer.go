package account

import (
	consumer2 "atlas-channel/kafka/consumer"
	"atlas-channel/server"
	"context"
	"github.com/Chronicle20/atlas-kafka/consumer"
	"github.com/Chronicle20/atlas-kafka/handler"
	"github.com/Chronicle20/atlas-kafka/message"
	"github.com/Chronicle20/atlas-kafka/topic"
	tenant "github.com/Chronicle20/atlas-tenant"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameAccountStatus = "account-status"
)

func StatusConsumer(l logrus.FieldLogger) func(groupId string) consumer.Config {
	return func(groupId string) consumer.Config {
		return consumer2.NewConfig(l)(consumerNameAccountStatus)(EnvEventTopicAccountStatus)(groupId)
	}
}

func StatusRegister(l logrus.FieldLogger, sc server.Model) (string, handler.Handler) {
	t, _ := topic.EnvProvider(l)(EnvEventTopicAccountStatus)()
	return t, message.AdaptHandler(message.PersistentConfig(handleAccountStatusEvent(sc)))
}

func handleAccountStatusEvent(sc server.Model) func(l logrus.FieldLogger, ctx context.Context, event statusEvent) {
	return func(l logrus.FieldLogger, ctx context.Context, event statusEvent) {
		t := sc.Tenant()
		if !t.Is(tenant.MustFromContext(ctx)) {
			return
		}

		if event.Status == EventAccountStatusLoggedIn {
			getRegistry().Login(Key{Tenant: t, Id: event.AccountId})
		} else if event.Status == EventAccountStatusLoggedOut {
			getRegistry().Logout(Key{Tenant: t, Id: event.AccountId})
		}
	}
}
