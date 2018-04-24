package main

import (
  "fmt"
  "net/http"
  "encoding/json"
  "github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
  "github.com/satori/go.uuid"
)

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

func GetRedisServer() bool {
  fmt.Println("Connecting to Redis server..")
  client = NewRedisServer()

  fmt.Println("PING")
  pong, err := client.Ping().Result()
  if err != nil {
    fmt.Println("Could not connect to Redis server:", err)
    return false
  }

  fmt.Println(pong)
  return true
}

// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/ping", pingHandler(formatter)).Methods("GET")
  mx.HandleFunc("/order", starbucksNewOrderHandler(formatter)).Methods("POST")
  mx.HandleFunc("/order/{id}", starbucksOrderStatusHandler(formatter)).Methods("GET")
  mx.HandleFunc("/orders", starbucksOrderStatusHandler(formatter)).Methods("GET")
}

// Starbucks API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Ping string }{"Pong"})
	}
}

// API Create New starbucks Order
func starbucksNewOrderHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
    var ord order
    err := json.NewDecoder(req.Body).Decode(&ord)
    if err != nil {
      fmt.Println(err)
      formatter.JSON(w, http.StatusBadRequest, err)
      return
    }

		uuid, _ := uuid.NewV4()
    ord.Id = uuid.String()
    ord.OrderStatus = "Order Placed"

		fmt.Println( "Order:", ord )
    key := ord.Id
    value, _ := json.Marshal(ord)
    err = client.Set(key, value, 0).Err()
    if err != nil {
      fmt.Println(err)
      formatter.JSON(w, http.StatusInternalServerError, err)
      return
  	}

		formatter.JSON(w, http.StatusOK, ord)
	}
}

// API Get Order Status
func starbucksOrderStatusHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		var uuid string = params["id"]
		if uuid == "" {
      keys := client.Keys("*")
      if keys == nil {
        fmt.Println("Order not found.")
        formatter.JSON(w, http.StatusNotFound, nil)
        return
      }
			var orders_array []string
			for key, value := range keys.Val() {
				fmt.Println("Key:", key, "Order:", value)
				orders_array = append(orders_array, value)
			}
			formatter.JSON(w, http.StatusOK, orders_array)
		} else {
			var ord, err = client.Get(uuid).Result()
      if err != nil {
        fmt.Println("Order not found.")
        formatter.JSON(w, http.StatusNotFound, err)
        return
      }
			fmt.Println("Order: ", ord)
			formatter.JSON(w, http.StatusOK, ord)
		}
	 }
 }
}
