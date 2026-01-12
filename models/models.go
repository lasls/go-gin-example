package models

import (
	"context"
	"log"
	"time"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database
var client *mongo.Client

type Model struct {
	ID         string    `bson:"_id,omitempty" json:"id"`
	CreatedOn  time.Time `bson:"created_on" json:"created_on"`
	ModifiedOn time.Time `bson:"modified_on" json:"modified_on"`
}

func init() {
	var (
		err      error
		connStr  string
		dbName   string
		username string
		password string
	)

	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Printf("Warning: Fail to get section 'database': %v", err)
		return
	}

	connStr = sec.Key("CONNECTION_STRING").String()
	dbName = sec.Key("database").String()
	username = sec.Key("username").String()
	password = sec.Key("password").String()

	// 如果提供了用户名和密码，则添加到连接字符串中
	if username != "" && password != "" {
		connStr = "mongodb://" + username + ":" + password + "@" + connStr[10:] // 去掉"mongodb://"前缀再重新拼接
	}

	// 设置客户端连接选项
	clientOptions := options.Client().ApplyURI(connStr)

	// 连接到MongoDB
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Printf("Warning: Failed to connect to MongoDB: %v", err)
		return
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Printf("Warning: Failed to ping MongoDB: %v", err)
		return
	}

	db = client.Database(dbName)

	log.Println("Connected to MongoDB!")

	// 初始化集合（相当于创建表）
	// 确保集合存在并创建索引
}

func CloseDB() {
	if client != nil {
		err := client.Disconnect(context.TODO())
		if err != nil {
			log.Println("Error disconnecting from MongoDB:", err)
		}
	}
	log.Println("Connection to MongoDB closed.")
}
