package main

import (
	"github.com/go-redis/redis"
)

type order struct {
	Id						string
	ProductName     string
	Category					string
	Price        	string
	
}

var client *redis.Client
