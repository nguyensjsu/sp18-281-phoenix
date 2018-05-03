package main

import (
  "github.com/go-redis/redis"
)

func NewRedisServer() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "13.57.96.105:6379",
		Password: "foobared",
		DB:       0,  // use default DB
	})

  return client
}
