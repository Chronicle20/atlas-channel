package session

import (
	"atlas-channel/kafka/producer"
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Destroy(l logrus.FieldLogger, kp producer.Provider) func(sessionId uuid.UUID, accountId uint32) {
	return func(sessionId uuid.UUID, accountId uint32) {
		l.Debugf("Destroying session for account [%d].", accountId)
		_ = kp(EnvCommandTopicAccountLogout)(logoutCommandProvider(sessionId, accountId))
	}
}

func UpdateState(l logrus.FieldLogger) func(ctx context.Context) func(sessionId uuid.UUID, accountId uint32, state int) (Model, error) {
	return func(ctx context.Context) func(sessionId uuid.UUID, accountId uint32, state int) (Model, error) {
		return func(sessionId uuid.UUID, accountId uint32, state int) (Model, error) {
			return updateState(l)(ctx)(sessionId, accountId, state)
		}
	}
}
