package handler

import (
	"atlas-channel/character"
	"atlas-channel/invite"
	"atlas-channel/party"
	"atlas-channel/session"
	"atlas-channel/socket/writer"
	"context"
	"github.com/Chronicle20/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

const (
	PartyOperationHandle       = "PartyOperationHandle"
	PartyOperationCreate       = "CREATE"
	PartyOperationLeave        = "LEAVE"
	PartyOperationExpel        = "EXPEL"
	PartyOperationChangeLeader = "CHANGE_LEADER"
	PartyOperationInvite       = "INVITE"
	PartyOperationJoin         = "JOIN"
)

func PartyOperationHandleFunc(l logrus.FieldLogger, ctx context.Context, wp writer.Producer) func(s session.Model, r *request.Reader, readerOptions map[string]interface{}) {
	return func(s session.Model, r *request.Reader, readerOptions map[string]interface{}) {
		op := r.ReadByte()
		if isPartyOperation(l)(readerOptions, op, PartyOperationCreate) {
			err := party.Create(l)(ctx)(s.CharacterId())
			if err != nil {
				l.WithError(err).Errorf("Character [%d] unable to attempt party creation.", s.CharacterId())
			}
			return
		}
		if isPartyOperation(l)(readerOptions, op, PartyOperationJoin) {
			partyId := r.ReadUint32()
			err := invite.Accept(l)(ctx)(s.CharacterId(), s.WorldId(), invite.InviteTypeParty, partyId)
			if err != nil {
				l.WithError(err).Errorf("Unable to issue invite acceptance command for character [%d].", s.CharacterId())
			}
			return
		}
		if isPartyOperation(l)(readerOptions, op, PartyOperationLeave) {
			p, err := party.GetByMemberId(l)(ctx)(s.CharacterId())
			if err != nil {
				l.WithError(err).Errorf("Unable to locate party for character [%d] to leave.", s.CharacterId())
				return
			}
			err = party.Leave(l)(ctx)(p.Id(), s.CharacterId())
			if err != nil {
				l.WithError(err).Errorf("Character [%d] unable to attempt leaving party.", s.CharacterId())
			}
			return
		}
		if isPartyOperation(l)(readerOptions, op, PartyOperationExpel) {
			targetCharacterId := r.ReadUint32()
			p, err := party.GetByMemberId(l)(ctx)(s.CharacterId())
			if err != nil {
				l.WithError(err).Errorf("Unable to locate party for character [%d] to leave.", s.CharacterId())
				return
			}
			err = party.Expel(l)(ctx)(p.Id(), s.CharacterId(), targetCharacterId)
			if err != nil {
				l.WithError(err).Errorf("Character [%d] unable to attempt expelling [%d] from party.", s.CharacterId(), targetCharacterId)
			}
			return
		}
		if isPartyOperation(l)(readerOptions, op, PartyOperationChangeLeader) {
			targetCharacterId := r.ReadUint32()
			p, err := party.GetByMemberId(l)(ctx)(s.CharacterId())
			if err != nil {
				l.WithError(err).Errorf("Unable to locate party for character [%d] to leave.", s.CharacterId())
				return
			}
			err = party.ChangeLeader(l)(ctx)(p.Id(), s.CharacterId(), targetCharacterId)
			if err != nil {
				l.WithError(err).Errorf("Character [%d] unable to pass leadership to [%d] in party.", s.CharacterId(), targetCharacterId)
			}
			return
		}
		if isPartyOperation(l)(readerOptions, op, PartyOperationInvite) {
			name := r.ReadAsciiString()
			cs, err := character.GetByName(l, ctx)(name)
			if err != nil {
				l.WithError(err).Errorf("Unable to locate character by name [%s] to invite to party.", name)
				err := session.Announce(l)(ctx)(wp)(writer.PartyOperation)(writer.PartyErrorBody(l)("UNABLE_TO_FIND_THE_CHARACTER", name))(s)
				if err != nil {
					return
				}
			}

			os, err := session.GetByCharacterId(s.Tenant(), s.WorldId(), s.ChannelId())(cs.Id())
			if err != nil || s.WorldId() != os.WorldId() || s.ChannelId() != os.ChannelId() {
				l.WithError(err).Errorf("Character [%d] not in channel. Cannot invite to party.", cs.Id())
				err = session.Announce(l)(ctx)(wp)(writer.PartyOperation)(writer.PartyErrorBody(l)("UNABLE_TO_FIND_THE_REQUESTED_CHARACTER_IN_THIS_CHANNEL", name))(s)
				if err != nil {
				}
				return
			}

			err = party.RequestInvite(l)(ctx)(s.CharacterId(), cs.Id())
			if err != nil {
				l.WithError(err).Errorf("Character [%d] was unable to request [%d] to join party.", s.CharacterId(), cs.Id())
			}
			return
		}
		l.Warnf("Character [%d] issued a unhandled party operation [%d].", s.CharacterId(), op)
	}
}

func isPartyOperation(l logrus.FieldLogger) func(options map[string]interface{}, op byte, key string) bool {
	return func(options map[string]interface{}, op byte, key string) bool {
		var genericCodes interface{}
		var ok bool
		if genericCodes, ok = options["operations"]; !ok {
			l.Errorf("Code [%s] not configured for use.", key)
			return false
		}

		var codes map[string]interface{}
		if codes, ok = genericCodes.(map[string]interface{}); !ok {
			l.Errorf("Code [%s] not configured for use.", key)
			return false
		}

		res, ok := codes[key].(float64)
		if !ok {
			l.Errorf("Code [%s] not configured for use.", key)
			return false
		}
		return byte(res) == op
	}
}
