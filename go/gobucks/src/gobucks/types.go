package main

import (
	"github.com/go-redis/redis"
)

type order struct {
	Id						string
	OrderName     string
	Size					string
	Options				string
	Price        	string
	OrderStatus 	string
}

var client *redis.Client
