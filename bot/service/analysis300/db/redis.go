package db

import (
	"eebot/g"
	"log"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     g.Config.GetString("analysis.redis.addr"),
		Password: g.Config.GetString("analysis.redis.password"),                                   // no password set
		DB:       g.Config.GetInt("analysis.redis.db"), // use default DB
	})
	if RDB == nil {
		log.Panic("panic: redis init")
	}
}
