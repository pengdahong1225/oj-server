package client

import (
	"oj-server/pkg/registry"
)

type Clients struct {
}

func NewClients(reg registry.Registry) (*Clients, error) {
	return &Clients{}, nil
}
