package handler

import (
	"atlas-channel/character/buff"
	"atlas-channel/session"
	"atlas-channel/socket/writer"
	"context"
	"github.com/Chronicle20/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

const CharacterBuffCancelHandle = "CharacterBuffCancel"

func CharacterBuffCancelHandleFunc(l logrus.FieldLogger, ctx context.Context, _ writer.Producer) func(s session.Model, r *request.Reader, readerOptions map[string]interface{}) {
	return func(s session.Model, r *request.Reader, readerOptions map[string]interface{}) {
		skillId := r.ReadUint32()
		_ = buff.Cancel(l)(ctx)(s.Map(), s.CharacterId(), skillId)
	}
}
