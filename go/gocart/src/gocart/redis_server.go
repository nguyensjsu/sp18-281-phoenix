package main

import (
  "github.com/go-redis/redis"
)

func NewRedisServer() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "13.56.224.231:8102",
		Password: "",
		DB:       0,  // use default DB
	})

  return client
}
