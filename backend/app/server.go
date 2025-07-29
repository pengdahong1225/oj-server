package app

import (
	"fmt"
	"oj-server/app/gateway"
	"oj-server/app/judge"
	"oj-server/app/problem"
	"oj-server/app/user"
	"oj-server/global"
	"oj-server/module/configManager"
)

type IServer interface {
	Init() error
	Run()
	Stop()
}

func NewServer() (IServer, error) {
	node := configManager.ServerConf

	var server IServer
	switch node.NodeType {
	case global.GatewayService:
		server = gateway.NewServer()
	case global.ProblemService:
		server = problem.NewServer()
	case global.JudgeService:
		server = judge.NewServer()
	case global.UserService:
		server = user.NewServer()
	default:
		return nil, fmt.Errorf("未知的节点类型: %s", node.NodeType)
	}
	return server, nil
}
