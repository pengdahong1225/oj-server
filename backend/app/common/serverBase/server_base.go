package ServerBase

import (
	"errors"
	"flag"
	"github.com/pengdahong1225/oj-server/backend/module/logger"
	"github.com/pengdahong1225/oj-server/backend/module/registry"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"github.com/pengdahong1225/oj-server/backend/module/utils"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
)

type Server struct {
	NodeType int
	NodeId   int
	Host     string
	Port     int
	Name     string
	SrvType  string // grpc | http
}

func (receiver *Server) Initialize() error {
	err := logger.Init("./log", "debug")
	if err != nil {
		return err
	}

	flag.IntVar(&receiver.NodeType, "node_type", 0, "node type")
	flag.IntVar(&receiver.NodeId, "node_id", 0, "node id")
	flag.StringVar(&receiver.Host, "host", "", "host")
	flag.IntVar(&receiver.Port, "port", 0, "port")
	flag.StringVar(&receiver.Name, "name", "", "name")
	flag.Parse()

	receiver.Host, err = utils.GetOutboundIPString()
	if err != nil {
		return err
	}

	logrus.Debugf("--------------- name:%v, node_type:%v, node_id:%v, host:%v, port:%v ---------------", receiver.Name, receiver.NodeType, receiver.NodeId, receiver.Host, receiver.Port)

	// 读取配置
	err = settings.Instance().LoadConfig()
	if err != nil {
		return err
	}

	return nil
}

// 服务注册
func (receiver *Server) Register() error {
	register, err := registry.NewRegistry()
	if err != nil {
		return err
	}

	if receiver.SrvType == "grpc" {
		err = register.RegisterServiceWithGrpc(&pb.PBNodeInfo{
			NodeType: int32(receiver.NodeType),
			NodeId:   int32(receiver.NodeId),
			Ip:       receiver.Host,
			Port:     int32(receiver.Port),
			State:    int32(pb.ENNodeState_EN_NODE_STATE_ONLINE),
			Name:     receiver.Name,
		})
		if err != nil {
			return err
		}
		return nil
	} else if receiver.SrvType == "http" {
		err = register.RegisterServiceWithHttp(&pb.PBNodeInfo{
			NodeType: int32(receiver.NodeType),
			NodeId:   int32(receiver.NodeId),
			Ip:       receiver.Host,
			Port:     int32(receiver.Port),
			State:    int32(pb.ENNodeState_EN_NODE_STATE_ONLINE),
			Name:     receiver.Name,
		})
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("unknown srv type")
	}
}

// 服务注销
func (receiver *Server) UnRegister() {
	register, err := registry.NewRegistry()
	if err != nil {
		logrus.Errorf("注销服务失败: %v", err)
		return
	}

	err = register.UnRegister(&pb.PBNodeInfo{
		NodeType: int32(receiver.NodeType),
		NodeId:   int32(receiver.NodeId),
		Ip:       receiver.Host,
		Port:     int32(receiver.Port),
		State:    int32(pb.ENNodeState_EN_NODE_STATE_OFFLINE),
	})
	if err != nil {
		logrus.Errorf("注销服务失败: %v", err)
		return
	}
}
