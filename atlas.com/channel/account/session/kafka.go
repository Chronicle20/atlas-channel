package session

import (
	"github.com/google/uuid"
)

const (
	EnvCommandTopic = "COMMAND_TOPIC_ACCOUNT_SESSION"

	CommandIssuerInternal = "INTERNAL"
	CommandIssuerLogin    = "LOGIN"
	CommandIssuerChannel  = "CHANNEL"

	CommandTypeProgressState = "PROGRESS_STATE"
	CommandTypeLogout        = "LOGOUT"
)

type command[E any] struct {
	SessionId uuid.UUID `json:"sessionId"`
	AccountId uint32    `json:"accountId"`
	Issuer    string    `json:"author"`
	Type      string    `json:"type"`
	Body      E         `json:"body"`
}

type progressStateCommandBody struct {
	State  uint8       `json:"state"`
	Params interface{} `json:"params"`
}

type logoutCommandBody struct {
}
