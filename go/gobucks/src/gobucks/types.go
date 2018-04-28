package main

import (
	"github.com/go-redis/redis"
)

type order struct {
	Id						string
	Items     []item
	TotalPrice        	string
	OrderStatus 	string
}

type item struct {
	Drink			string
	Size					string
	Options				string
	Price        	string
	Quantity     int
}

var client *redis.Client
