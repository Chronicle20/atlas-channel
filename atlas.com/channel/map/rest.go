package _map

import "strconv"

type RestModel struct {
	Id string `json:"-"`
}

func Extract(m RestModel) (uint32, error) {
	id, err := strconv.ParseUint(m.Id, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}