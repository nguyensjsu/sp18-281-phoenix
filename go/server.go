package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"github.com/codegangsta/negroni"
	"github.com/streadway/amqp"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

// MongoDB Config
var mongodb_server = "mongodb"
var mongodb_database = "cmpe281"
var mongodb_collection = "starbucks"



// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	n.UseHandler(mx)
	return n
}

// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/ping", pingHandler(formatter)).Methods("GET")
	mx.HandleFunc("/starbucks", starbucksUpdateHandler(formatter)).Methods("PUT")
    mx.HandleFunc("/order", starbucksNewOrderHandler(formatter)).Methods("POST")
}

// Helper Functions
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// Starbucks API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{" Starbucks API version 1.0 alive!"})
	}
}


// API Create New starbucks Order
func starbucksNewOrderHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		uuid := uuid.NewV4()
    	var ord = order {
					Id: uuid.String(),            		
					OrderStatus: "Order Placed",
		}
		if orders == nil {
			orders = make(map[string]order)
		}
		orders[uuid.String()] = ord
		queue_send(uuid.String())
		fmt.Println( "Orders: ", orders )
		formatter.JSON(w, http.StatusOK, ord)
	}
}


