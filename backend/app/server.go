package app

import (
	"fmt"
	"oj-server/app/gateway"
	"oj-server/global"
	"oj-server/module/configManager"
)

type IServer interface {
	Init() error
	Run()
}

func NewServer() (IServer, error) {
	node := configManager.ServerConf

	var server IServer
	switch node.NodeType {
	case global.GatewayService:
		server = gateway.NewServer()
	default:
		return nil, fmt.Errorf("未知的节点类型: %s", node.NodeType)
	}
	return server, nil
}
