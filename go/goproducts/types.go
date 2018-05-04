package main

import (
	"github.com/go-redis/redis"
)

type product struct {
	Id              string
	ProductName     string
	Category		string
	Price        	string
	Rating          int
	
}

var client *redis.Client


