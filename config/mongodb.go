package config

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

//客户端
var Client *mongo.Client

//操作数据库
var Collection *mongo.Collection

func Initmongo(addr string) {
	clientOptions := options.Client().ApplyURI(addr)

	// 连接到MongoDB
	var err error
	Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		logrus.Fatalf("连接 config 服务器失败: %v", err)
	}
	// 检查连接
	err = Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	Collection = Client.Database("alms1").Collection("user_info")

}
