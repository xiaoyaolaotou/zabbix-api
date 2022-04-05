package models

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"log"
)

var RdsClient *redis.Client

func RedisInit()error {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	_,err:=client.Ping().Result()
	if err != nil {
		log.Fatal("Redis连接失败")
		return err
	}
	RdsClient = client
	return nil
}
