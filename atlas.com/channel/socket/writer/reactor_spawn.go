package writer

import (
	"atlas-channel/reactor"
	"github.com/Chronicle20/atlas-socket/response"
	"github.com/Chronicle20/atlas-tenant"
	"github.com/sirupsen/logrus"
)

const (
	ReactorSpawn = "ReactorSpawn"
)

func ReactorSpawnBody(l logrus.FieldLogger, t tenant.Model) func(m reactor.Model) BodyProducer {
	return func(m reactor.Model) BodyProducer {
		return func(w *response.Writer, options map[string]interface{}) []byte {
			w.WriteInt(m.Id())
			w.WriteInt(m.Classification())
			w.WriteInt8(m.State())
			w.WriteInt16(m.X())
			w.WriteInt16(m.Y())
			w.WriteByte(m.Direction())
			w.WriteAsciiString(m.Name())
			return w.Bytes()
		}
	}
}
