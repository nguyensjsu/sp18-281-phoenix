package main



type order struct {
	Id             	string 	
	OrderStatus 	string	
}

var orders map[string] order
