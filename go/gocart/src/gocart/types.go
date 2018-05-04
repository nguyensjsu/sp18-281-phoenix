package main

import (
	"github.com/go-redis/redis"
)

type Cart struct {
	Id						string
	Items     []Item
	TotalPrice        	float64
	Status 	string
}

type Item struct {
	Drink			string
	Size					string
	Options				string
	Price        	float64
	Quantity     int
}

var client *redis.Client
