package handler

import (
	"atlas-channel/character/expression"
	"atlas-channel/session"
	"atlas-channel/socket/writer"
	"context"
	"github.com/Chronicle20/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

const CharacterExpressionHandle = "CharacterExpressionHandle"

func CharacterExpressionHandleFunc(l logrus.FieldLogger, ctx context.Context, _ writer.Producer) func(s session.Model, r *request.Reader, readerOptions map[string]interface{}) {
	return func(s session.Model, r *request.Reader, readerOptions map[string]interface{}) {
		emote := r.ReadUint32()
		_ = expression.Change(l)(ctx)(s.CharacterId(), s.Map(), emote)
	}
}
