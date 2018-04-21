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
  	mx.HandleFunc("/order", starbucksNewOrderHandler(formatter)).Methods("PUT")
	mx.HandleFunc("/order/{id}", starbucksOrderStatusHandler(formatter)).Methods("GET")
	mx.HandleFunc("/order/{id}", starbucksUpdateOrderHandler(formatter)).Methods("POST")
	mx.HandleFunc("/order", starbucksOrderStatusHandler(formatter)).Methods("GET")
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

// API Update starbucks Order
func starbucksUpdateOrderHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
    	var m order
    	_ = json.NewDecoder(req.Body).Decode(&m)
    	fmt.Println("Update Order To: ", m.CountGumballs)
		session, err := mgo.Dial(mongodb_server)
        if err != nil {
                panic(err)
        }
        defer session.Close()
        session.SetMode(mgo.Monotonic, true)
        c := session.DB(mongodb_database).C(mongodb_collection)
        query := bson.M{"Id" : m.id}
        change := bson.M{"$set": bson.M{ "order" : m.order}}
        err = c.Update(query, change)
        if err != nil {
                log.Fatal(err)
        }
       	var result bson.M
        err = c.Find(bson.M{"Id" : m.id}).One(&result)
        if err != nil {
                log.Fatal(err)
        }
        fmt.Println("Order:", result )
		formatter.JSON(w, http.StatusOK, result)
	}
}

// API Get Order Status
func starbucksOrderStatusHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		var uuid string = params["id"]
		fmt.Println("Order ID: ", uuid)
		if uuid == "" {
			fmt.Println("Orders:", orders)
			var orders_array []order
			for key, value := range orders {
				fmt.Println("Key:", key, "Value:", value)
				orders_array = append(orders_array, value)
			}
			formatter.JSON(w, http.StatusOK, orders_array)
		} else {
			var ord = orders[uuid]
			fmt.Println("Order: ", ord)
			formatter.JSON(w, http.StatusOK, ord)
		}
	}
}

// API Process Orders
func starbucksProcessOrdersHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		for key, value := range orders {
			fmt.Println("Key:", key, "Value:", value)
			var ord = orders[key]
			ord.OrderStatus = "Order Processed"
			orders[key] = ord
		}
		fmt.Println("Orders: ", orders)
		formatter.JSON(w, http.StatusOK, "Orders Processed!")
	}
}
