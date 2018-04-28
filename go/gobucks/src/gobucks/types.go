package main

import (
	"github.com/go-redis/redis"
)

type Order struct {
	Id						string
	Items     []Item
	TotalPrice        	string
	OrderStatus 	string
}

type Item struct {
	Drink			string
	Size					string
	Options				string
	Price        	string
	Quantity     int
}

var client *redis.Client
