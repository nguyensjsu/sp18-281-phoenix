package main

import (
	"fmt"
  //"log"
  "net/http"
  "encoding/json"
	"github.com/gorilla/mux"
  "github.com/satori/go.uuid"
  "github.com/unrolled/render"
  "github.com/codegangsta/negroni"
)

// Init Background Processes
func init() {

}

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
	mx.HandleFunc("/cart", newCartHandler(formatter)).Methods("POST")
	mx.HandleFunc("/cart/{id}", getCartHandler(formatter)).Methods("GET")
	mx.HandleFunc("/cart/{id}", updateCartHandler(formatter)).Methods("PUT")
	mx.HandleFunc("/cart/{id}", removeCartHandler(formatter)).Methods("DELETE")
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"API works!"})
	}
}

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

    totalPrice := float64(0)
    for _, num := range crt.Items {
      totalPrice += float64(num.Price * float64(num.Quantity))
    }
    crt.TotalPrice = totalPrice

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

func updateCartHandler(formatter *render.Render) http.HandlerFunc {
  return func(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    var uuid string = params["id"]
    var result, err = client.Get(uuid).Result()
    if err != nil {
      fmt.Println("Cart not found")
      formatter.JSON(w, http.StatusOK, err)
      return
    }
    var crt Cart
    json.Unmarshal([]byte(result), &crt)

    var cart Cart
    err = json.NewDecoder(req.Body).Decode(&cart)
    if err != nil {
      fmt.Println(err)
      formatter.JSON(w, http.StatusBadRequest, err)
      return
    }

    crt = cart
    crt.Id = uuid
    crt.Status = "Cart Updated"
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