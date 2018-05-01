package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"net/http"
	"log"
	"github.com/go-redis/redis"
)

// Init Background Processes
func init() {

}

func NewRedisServer() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "13.56.224.231:8102",
		Password: "",
		DB:       0,  // use default DB
	})

  return client
}

// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/ping", pingHandler(formatter)).Methods("GET")
	mx.HandleFunc("/cart", newCartHandler(formatter)).Methods("PUT")
	mx.HandleFunc("/cart/{id}", getCartHandler(formatter)).Methods("GET")
	mx.HandleFunc("/cart/{id}", UpdateCartHandler(formatter)).Methods("POST")
	mx.HandleFunc("/cart/{id}", removeCartHandler(formatter)).Methods("DELETE")
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"API works!"})
	}
}

// API redis Handler
/*func redHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		conn, err := redis.Dial("tcp", "ec2-54-183-158-59.us-west-1.compute.amazonaws.com:8102")
   	client = NewRedisServer()

	err = client.Set(params["id"], params["value"], 0).Err()
    if err != nil {
      fmt.Println(err)
      formatter.JSON(w, http.StatusInternalServerError, err)
      return
  	}

		formatter.JSON(w, http.StatusOK, "key inserted.")
	}
}*/

// API Create New starbucks Cart
func newCartHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
    var crt Cart
    err := json.NewDecoder(req.Body).Decode(&crt)
    if err != nil {
      fmt.Println(err)
      formatter.JSON(w, http.StatusBadRequest, err)
      return
    }

	uuid, _ := uuid.NewV4()
    crt.Id = uuid.String()
    crt.Status = "Added New Cart"

	fmt.Println( "Cart:", crt )
    key := crt.Id
    value, _ := json.Marshal(crt)
    err = client.Set(key, value, 0).Err()
    if err != nil {
      fmt.Println(err)
      formatter.JSON(w, http.StatusInternalServerError, err)
      return
  	}

		formatter.JSON(w, http.StatusOK, crt)
	}
}

// API Get Order Status
func getCartHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		var uuid string = params["id"]
		var crt, err = client.Get(uuid).Result()
	      if err != nil {
	        fmt.Println("Cart not found.")
	        formatter.JSON(w, http.StatusNotFound, err)
	        return
	      }
			fmt.Println("Cart: ", crt)
      var cart Cart
      json.Unmarshal([]byte(crt), &cart)
			formatter.JSON(w, http.StatusOK, cart)
		}
	}
}

// API delete cart
func removeCartHandler(formatter *render.Render) http.HandlerFunc {
  return func(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    var uuid string = params["id"]
    var crt, err = client.Get(uuid).Result()
    if err != nil {
      fmt.Println("Cart not found.")
      formatter.JSON(w, http.StatusNotFound, err)
      return
    }
    var cart Cart
    json.Unmarshal([]byte(crt), &cart)
    err = client.Del(uuid).Err()
    if err != nil {
      fmt.Println(err)
      formatter.JSON(w, http.StatusInternalServerError, err)
      return
    }
    cart.Status = "Removed Cart"
    formatter.JSON(w, http.StatusOK, cart)
	}
}