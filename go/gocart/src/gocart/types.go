package main

import (
	"github.com/go-redis/redis"
)

type Cart struct {
	Id						string
	Items     []Item
	TotalPrice        	string
	Status 	string
}

type Item struct {
	Drink			string
	Size					string
	Options				string
	Price        	string
	Quantity     int
}

var client *redis.Client
