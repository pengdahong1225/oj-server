package global

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var (
	prefix            = "/app/log/web-service"
	ConfigInstance    config
	RedisPoolInstance *redis.Pool
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

// 初始化日志信息
func initLog() error {
	var path = ConfigInstance.System_.Prefix + ".log." + time.Now().Format("2006-01-02")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Llongfile)
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
