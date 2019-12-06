package Services

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

//--创建空间--
type SpaceCreateRequest struct {
	Token    string `json:"token"`
	Level    int    `json:"level"`
	Province string `json:"province"`
	City     string `json:"city"`
	Area     string `json:"area"`
	District string `json:"distinct"`
	Building string `json:"building"`
	Storey   string `json:"storey"`
	Room     string `json:"room"`
	Place    string `json:"place"`
	Title    string `json:"title"`
}

func SpaceCreateEndpoint(spacecreateservice WSpaceCreateService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*SpaceCreateRequest)
		result := spacecreateservice.NewSpace(r)
		return result, nil
	}
}
