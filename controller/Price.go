package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var RedisClient *redis.Client

func GetPrice(w http.ResponseWriter, r *http.Request) {
	val, err := RedisClient.Get(ctx, "TV").Result()
	if err != nil {
		log.Println(err)
	}
	w.Write([]byte("TV: " + val))
}

func SetPrice(w http.ResponseWriter, r *http.Request) {
	err := RedisClient.Set(ctx, "TV", "6000", time.Minute).Err()
	if err != nil {
		log.Println(err)
	}
}
