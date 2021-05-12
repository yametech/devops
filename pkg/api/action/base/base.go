package base

import (
	"fmt"
	"github.com/yametech/devops/pkg/api"
)

type baseServer struct {
	*api.Server
}

func NewBaseServer(serviceName string, server *api.Server) *baseServer {
	base := &baseServer{
		Server: server,
	}
	group := base.Group(fmt.Sprintf("/%s", serviceName))

	_ = group

	return base
}
