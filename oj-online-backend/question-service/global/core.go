package global

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"log"
	"os"
)

var (
	ConfigInstance    config
	RedisPoolInstance *redis.Pool
	MqConnection      *amqp.Connection
)

// 1.读取配置
// 2.初始化日志
// 3.redis
func init() {
	if err := loadConfig(); err != nil {
		panic(err)
	}
	if err := initLog(); err != nil {
		panic(err)
	}
	if err := initRedis(); err != nil {
		panic(err)
	}
	if err := initMQ(); err != nil {
		panic(err)
	}
}

func loadConfig() error {
	// 读取环境变量
	logMode := os.Getenv("LOG_MODE")
	// 读取配置
	viperConfig := viper.New()
	if logMode == "release" {
		viperConfig.SetConfigName("prod")
	} else {
		viperConfig.SetConfigName("dev")
	}
	viperConfig.AddConfigPath("config")
	viperConfig.SetConfigType("yaml")
	if e := viperConfig.ReadInConfig(); e != nil {
		return e
	}
	if e := viperConfig.Unmarshal(&ConfigInstance); e != nil {
		return e
	}
	// 配置热重载
	viper.OnConfigChange(func(event fsnotify.Event) {
		log.Println("config file changed:", event.Name)
		if e := viperConfig.Unmarshal(&ConfigInstance); e != nil {
			log.Println("config file update failed:", event.Name)
		}
	})
	// 监听配置文件
	viper.WatchConfig()
	return nil
}

func initRedis() error {
	dsn := fmt.Sprintf("%s:%d", ConfigInstance.Redis_.Ip, ConfigInstance.Redis_.Port)
	RedisPoolInstance = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", dsn)
		},
		DialContext:     nil,
		TestOnBorrow:    nil,
		MaxIdle:         8,   // 最大空闲数
		MaxActive:       10,  // 最大连接数 -- 0表示无限制
		IdleTimeout:     100, // 最大空闲时间
		Wait:            false,
		MaxConnLifetime: 0,
	}
	return nil
}

func initMQ() error {
	var err error
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		ConfigInstance.Mq_.User,
		ConfigInstance.Mq_.Pass,
		ConfigInstance.Mq_.Host,
		ConfigInstance.Mq_.Port,
		ConfigInstance.Mq_.VHost,
	)
	MqConnection, err = amqp.Dial(dsn)
	return err
}
