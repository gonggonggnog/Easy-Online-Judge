package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var DB = Init()

var ctx = context.Background()

var RDB = RdbInit()

func RedisSet(key, value string) {
	RDB.Set(ctx, key, value, time.Minute*10)
}

func RedisGet(key string) (string, error) {
	return RDB.Get(ctx, "name").Result()
}
func Init() *gorm.DB {
	dsn := "root:gong123123@tcp(47.115.224.170:8988)/ginoj?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("连接失败")
		fmt.Println(err)

	}
	return db
}
func RdbInit() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "47.115.224.170:6379",
		Password: "",
		DB:       0,
	})
}
