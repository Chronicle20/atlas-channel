package party

import (
	"atlas-channel/rest"
	"fmt"
	"github.com/Chronicle20/atlas-rest/requests"
)

const (
	Resource        = "parties"
	ByMemberId      = Resource + "?filter[members.id]=%d"
	ById            = Resource + "/%d"
	MembersResource = ById + "/members"
)

func getBaseRequest() string {
	return requests.RootUrl("PARTIES")
}

func requestById(id uint32) requests.Request[RestModel] {
	return rest.MakeGetRequest[RestModel](fmt.Sprintf(getBaseRequest()+ById, id))
}

func requestByMemberId(id uint32) requests.Request[[]RestModel] {
	return rest.MakeGetRequest[[]RestModel](fmt.Sprintf(getBaseRequest()+ByMemberId, id))
}

func requestMembers(id uint32) requests.Request[[]MemberRestModel] {
	return rest.MakeGetRequest[[]MemberRestModel](fmt.Sprintf(getBaseRequest()+MembersResource, id))
}
