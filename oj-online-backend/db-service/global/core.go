package global

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gomodule/redigo/redis"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

// 初始化全局DB连接
var (
	err               error
	ConfigInstance    config
	DBInstance        *gorm.DB
	RedisPoolInstance *redis.Pool
)

// 1.读取配置
// 2.初始化日志
// 3.db连接
// 4.redis
func init() {
	if err = loadConfig(); err != nil {
		panic(err)
	}
	if err = initLog(); err != nil {
		panic(err)
	}
	if err = initDB(); err != nil {
		panic(err)
	}
	if err = initRedis(); err != nil {
		panic(err)
	}
	AntsPoolInstance, _ = ants.NewPool(ants.DefaultAntsPoolSize, ants.WithPanicHandler(AntsPanicHandler))
}

func loadConfig() error {
	// 读取环境变量
	logMode := os.Getenv("LOG_MODE")
	// 读取配置文件
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

func initDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", ConfigInstance.Sql_.User,
		ConfigInstance.Sql_.Pwd, ConfigInstance.Sql_.Host, ConfigInstance.Sql_.Port, ConfigInstance.Sql_.Db)
	var e error
	DBInstance, e = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// SkipDefaultTransaction: true, //全局禁用默认事务
	})
	return e
}

func initRedis() error {
	dsn := fmt.Sprintf("%s:%d", ConfigInstance.Redis_.Host, ConfigInstance.Redis_.Port)
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
