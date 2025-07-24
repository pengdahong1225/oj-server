package ServerBase

import (
	"flag"
	"github.com/sirupsen/logrus"
	"oj-server/module/logger"
	"oj-server/module/registry"
	"oj-server/module/settings"
	"oj-server/proto/pb"
	"oj-server/utils"
)

type Server struct {
	NodeType int
	NodeId   int
	Host     string
	Port     int
	Name     string
	Scheme   string // grpc | http
	Config   string
}

func (receiver *Server) Initialize() error {
	flag.IntVar(&receiver.NodeType, "node_type", 0, "node type")
	flag.IntVar(&receiver.NodeId, "node_id", 0, "node id")
	flag.StringVar(&receiver.Host, "host", "", "host")
	flag.IntVar(&receiver.Port, "port", 0, "port")
	flag.StringVar(&receiver.Name, "name", "", "name")
	flag.StringVar(&receiver.Config, "config", "", "config")
	flag.Parse()

	err := logger.Init("./log", receiver.Name, logrus.DebugLevel)
	if err != nil {
		return err
	}

	receiver.Host, err = utils.GetOutboundIPString()
	if err != nil {
		return err
	}

	logrus.Debugf("--------------- name:%v, node_type:%v, node_id:%v, host:%v, port:%v ---------------", receiver.Name, receiver.NodeType, receiver.NodeId, receiver.Host, receiver.Port)

	// 读取配置
	err = settings.AppConf.LoadConfig(receiver.Config)
	if err != nil {
		return err
	}

	// 初始化注册中心
	err = registry.Init(receiver.Scheme)
	if err != nil {
		return err
	}

	return nil
}

// 服务注册
func (receiver *Server) Register() error {
	err := registry.RegisterService(&pb.PBNodeInfo{
		NodeType: pb.ENNodeType(receiver.NodeType),
		NodeId:   int32(receiver.NodeId),
		Ip:       receiver.Host,
		Port:     int32(receiver.Port),
		Name:     receiver.Name,
	})
	if err != nil {
		return err
	}
	return nil
}

// 服务注销
func (receiver *Server) UnRegister() {
	err := registry.DeregisterService(&pb.PBNodeInfo{
		NodeType: pb.ENNodeType(receiver.NodeType),
		NodeId:   int32(receiver.NodeId),
		Ip:       receiver.Host,
		Port:     int32(receiver.Port),
	})
	if err != nil {
		logrus.Errorf("注销服务失败: %v", err)
		return
	}
}
