package config

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

// RDB 是一个 config 服务器的客户端实例
var RDB *redis.Client

// InitRedis 负责初始化 config 客户端
func InitRedis(addr string) {
	// 连接 config 服务器
	RDB = redis.NewClient(&redis.Options{Addr: addr})
	// 测试是否连接成功
	_, err := RDB.Ping().Result()
	// 如果连接失败直接退出程序
	if err != nil {
		logrus.Fatalf("连接 config 服务器失败: %v", err)
	}
}
