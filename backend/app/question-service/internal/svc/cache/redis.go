package cache

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
)

var (
	Rdb  *redis.Client
	once sync.Once
)

func Init(ip string, port int) {
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%d", ip, port)
		Rdb = redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr:    dsn,
		})
	})
}
