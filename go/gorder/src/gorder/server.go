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
	mx.HandleFunc("/", pingHandler(formatter)).Methods("GET")
  mx.HandleFunc("/order", starbucksCreateOrderHandler(formatter)).Methods("POST")
  mx.HandleFunc("/order/{id}", starbucksReadOrderHandler(formatter)).Methods("GET")
  mx.HandleFunc("/orders", starbucksReadOrderHandler(formatter)).Methods("GET")
  mx.HandleFunc("/order/{id}", starbucksUpdateHandler(formatter)).Methods("PUT")
  mx.HandleFunc("/order/{id}", starbucksUpdateHandler(formatter)).Methods("PATCH")
  mx.HandleFunc("/order/{id}", starbucksDeleteOrderHandler(formatter)).Methods("DELETE")
}

// Starbucks API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Ping string }{"Pong"})
	}
}

// API Create Order
func starbucksCreateOrderHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
    var ord Order
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

// API Read Order(s)
func starbucksReadOrderHandler(formatter *render.Render) http.HandlerFunc {
  return func(w http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		var uuid string = params["id"]
		if uuid == "" {
      keys := client.Keys("*")
      if keys == nil {
        fmt.Println("Order not found")
        formatter.JSON(w, http.StatusOK, nil)
        return
      }
			var orders []Order
			for _, value := range keys.Val() {
        var order Order
        var ord, err = client.Get(value).Result()
        if err != nil {
          continue
        }
        fmt.Println("Order:", ord)
        json.Unmarshal([]byte(ord), &order)
				orders = append(orders, order)
			}
			formatter.JSON(w, http.StatusOK, orders)
		} else {
      var order Order
			var ord, err = client.Get(uuid).Result()
      if err != nil {
        fmt.Println("Order not found")
        formatter.JSON(w, http.StatusOK, nil)
        return
      }
			fmt.Println("Order:", ord)
      json.Unmarshal([]byte(ord), &order)
			formatter.JSON(w, http.StatusOK, order)
		}
	}
}

// API Update Order
func starbucksUpdateHandler(formatter *render.Render) http.HandlerFunc {
  return func(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    var uuid string = params["id"]
    var result, err = client.Get(uuid).Result()
    if err != nil {
      fmt.Println("Order not found")
      formatter.JSON(w, http.StatusOK, err)
      return
    }
    var ord Order
    json.Unmarshal([]byte(result), &ord)

    var order Order
    err = json.NewDecoder(req.Body).Decode(&order)
    if err != nil {
      fmt.Println(err)
      formatter.JSON(w, http.StatusBadRequest, err)
      return
    }

    switch req.Method {
    case "PATCH":
      ord.OrderStatus = order.OrderStatus
    case "PUT":
      ord = order
      ord.Id = uuid
      ord.OrderStatus = "Order Placed"
    }
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

// API Delete Order
func starbucksDeleteOrderHandler(formatter *render.Render) http.HandlerFunc {
  return func(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    var uuid string = params["id"]
    var _, err = client.Get(uuid).Result()
    if err != nil {
      fmt.Println("Order not found")
      formatter.JSON(w, http.StatusOK, err)
      return
    }

    err = client.Del(uuid).Err()
    if err != nil {
      fmt.Println(err)
      formatter.JSON(w, http.StatusInternalServerError, err)
      return
  	}
    formatter.JSON(w, http.StatusOK, "OK")
  }
}
