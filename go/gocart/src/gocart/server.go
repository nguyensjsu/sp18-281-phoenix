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

// Init Background Processes
func init() {

}

func NewRedisServer() *redis.Client {
	client := redis.NewClient(&edis.Options{
		Addr:     "54.67.22.208:8102",
		Password: "",
		DB:       0,  // use default DB
	})

  return client
}

// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/ping", pingHandler(formatter)).Methods("GET")
	mx.HandleFunc("/redis", redHandler(formatter)).Methods("GET")
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"API works!"})
	}
}

// API redis Handler
func redHandler(formatter *render.Render) http.HandlerFunc {
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
}