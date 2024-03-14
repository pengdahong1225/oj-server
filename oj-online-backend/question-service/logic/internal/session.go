package internal

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type Session struct {
	Id            string
	WebConnection *websocket.Conn
	Timer         *time.Timer
	Msg           []byte
}

var (
	once       sync.Once
	SessionMap map[string]*Session
)

func NewSession(userId int64) *Session {
	// 延迟初始化
	once.Do(func() {
		if SessionMap == nil {
			SessionMap = make(map[string]*Session, 10)
		}
	})

	sessionID := fmt.Sprintf("%d-%d", userId, time.Now().Unix())
	s := &Session{
		Id:            sessionID,
		WebConnection: nil,
		Timer:         nil,
		Msg:           nil,
	}
	SessionMap[sessionID] = s
	return s
}
