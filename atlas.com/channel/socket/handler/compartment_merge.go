package handler

import (
	"atlas-channel/compartment"
	"atlas-channel/session"
	"atlas-channel/socket/writer"
	"context"
	"github.com/Chronicle20/atlas-constants/inventory"
	"github.com/Chronicle20/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

const (
	CompartmentMerge = "CompartmentMerge"
)

func CompartmentMergeHandleFunc(l logrus.FieldLogger, ctx context.Context, _ writer.Producer) func(s session.Model, r *request.Reader, readerOptions map[string]interface{}) {
	return func(s session.Model, r *request.Reader, readerOptions map[string]interface{}) {
		updateTime := r.ReadUint32()
		compartmentType := inventory.Type(r.ReadByte())

		isValid := false
		for _, validType := range inventory.Types {
			if compartmentType == validType {
				isValid = true
				break
			}
		}

		if !isValid {
			l.Warnf("Character [%d] issued compartment merge with invalid compartment type [%d].", s.CharacterId(), compartmentType)
			return
		}

		err := compartment.NewProcessor(l, ctx).Merge(s.CharacterId(), compartmentType, updateTime)
		if err != nil {
			l.WithError(err).Errorf("Failed to send compartment merge command for character [%d].", s.CharacterId())
		}
	}
}