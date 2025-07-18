package messenger

import (
	"github.com/Chronicle20/atlas-constants/channel"
	"github.com/Chronicle20/atlas-constants/world"
	"github.com/jtumidanski/api2go/jsonapi"
	"strconv"
)

type RestModel struct {
	Id      uint32            `json:"-"`
	Members []MemberRestModel `json:"-"`
}

func (r RestModel) GetName() string {
	return "messengers"
}

func (r RestModel) GetID() string {
	return strconv.Itoa(int(r.Id))
}

func (r *RestModel) SetID(strId string) error {
	id, err := strconv.Atoi(strId)
	if err != nil {
		return err
	}
	r.Id = uint32(id)
	return nil
}

func (r RestModel) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "members",
			Name: "members",
		},
	}
}

func (r RestModel) GetReferencedIDs() []jsonapi.ReferenceID {
	var result []jsonapi.ReferenceID
	for _, v := range r.Members {
		result = append(result, jsonapi.ReferenceID{
			ID:   v.GetID(),
			Type: "members",
			Name: "members",
		})
	}
	return result
}

func (r RestModel) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	var result []jsonapi.MarshalIdentifier
	for key := range r.Members {
		result = append(result, r.Members[key])
	}

	return result
}

func (r *RestModel) SetToManyReferenceIDs(name string, IDs []string) error {
	if name == "members" {
		for _, ID := range IDs {
			id, err := strconv.Atoi(ID)
			if err != nil {
				return err
			}
			r.Members = append(r.Members, MemberRestModel{
				Id:        uint32(id),
				Name:      "",
				WorldId:   0,
				ChannelId: 0,
				Online:    false,
				Slot:      0,
			})
		}
	}
	return nil
}

func (r *RestModel) SetReferencedStructs(references map[string]map[string]jsonapi.Data) error {
	if refMap, ok := references["members"]; ok {
		var nm []MemberRestModel
		for _, m := range r.Members {
			if data, ok := refMap[m.GetID()]; ok {
				srm := MemberRestModel{}
				err := jsonapi.ProcessIncludeData(&srm, data, references)
				if err != nil {
					return err
				}
				err = srm.SetID(m.GetID())
				if err != nil {
					return err
				}
				nm = append(nm, srm)
			}
		}
		r.Members = nm
	}
	return nil
}

func Extract(rm RestModel) (Model, error) {
	var members = make([]MemberModel, 0)
	for _, m := range rm.Members {
		mm, err := ExtractMember(m)
		if err != nil {
			return Model{}, err
		}
		members = append(members, mm)
	}

	return Model{
		id:      rm.Id,
		members: members,
	}, nil
}

func ExtractMember(rm MemberRestModel) (MemberModel, error) {
	return MemberModel{
		id:        rm.Id,
		name:      rm.Name,
		worldId:   rm.WorldId,
		channelId: rm.ChannelId,
		online:    rm.Online,
		slot:      rm.Slot,
	}, nil
}

type MemberRestModel struct {
	Id        uint32 `json:"-"`
	Name      string `json:"name"`
	WorldId   world.Id   `json:"worldId"`
	ChannelId channel.Id `json:"channelId"`
	Online    bool   `json:"online"`
	Slot      byte   `json:"slot"`
}

func (r MemberRestModel) GetName() string {
	return "members"
}

func (r MemberRestModel) GetID() string {
	return strconv.Itoa(int(r.Id))
}

func (r *MemberRestModel) SetID(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	r.Id = uint32(id)
	return nil
}
