package client

import (
	"context"
	"fmt"
	"oj-server/common/consts"
	"oj-server/common/proto/pb"
	"oj-server/pkg/registry"
)

type Clients struct {
	User pb.UserServiceClient
}

func NewClients(reg registry.Registry) (*Clients, error) {
	ctx := context.Background()

	userConn, err := reg.Dial(ctx, consts.UserService)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", consts.UserService, err)
	}

	return &Clients{
		User: pb.NewUserServiceClient(userConn),
	}, nil
}
