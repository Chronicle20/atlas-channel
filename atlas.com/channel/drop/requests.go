package drop

import (
	"atlas-channel/rest"
	"fmt"
	"github.com/Chronicle20/atlas-rest/requests"
)

const (
	Resource = "worlds/%d/channels/%d/maps/%d/drops"
)

func getBaseRequest() string {
	return requests.RootUrl("DROPS")
}

func requestInMap(worldId byte, channelId byte, mapId uint32) requests.Request[[]RestModel] {
	return rest.MakeGetRequest[[]RestModel](fmt.Sprintf(getBaseRequest()+Resource, worldId, channelId, mapId))
}
