package map_

import (
	"atlas-channel/rest"
	"fmt"
	_map "github.com/Chronicle20/atlas-constants/map"
	"github.com/Chronicle20/atlas-rest/requests"
)

const (
	getMap = "data/maps/%d"
)

func getBaseRequest() string {
	return requests.RootUrl("DATA")
}

func requestMap(mapId _map.Id) requests.Request[RestModel] {
	return rest.MakeGetRequest[RestModel](fmt.Sprintf(getBaseRequest()+getMap, mapId))
}
