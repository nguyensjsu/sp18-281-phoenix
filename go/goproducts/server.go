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
  mx.HandleFunc("/product", starbucksCreateProductHandler(formatter)).Methods("POST")
  mx.HandleFunc("/product/{id}", starbucksProductDetailsHandler(formatter)).Methods("GET")
  mx.HandleFunc("/products", starbucksProductDetailsHandler(formatter)).Methods("GET")
  mx.HandleFunc("/product/{id}", starbucksDeleteProductHandler(formatter)).Methods("DELETE")
}

// Starbucks API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
  return func(w http.ResponseWriter, req *http.Request) {
    formatter.JSON(w, http.StatusOK, struct{ Ping string }{"Pong"})
  }
}

// API Create Order
func starbucksCreateProductHandler(formatter *render.Render) http.HandlerFunc {
  return func(w http.ResponseWriter, req *http.Request) {
    var prd product
    err := json.NewDecoder(req.Body).Decode(&prd)
    if err != nil {
      fmt.Println(err)
      formatter.JSON(w, http.StatusBadRequest, err)
      return
    }

    uuid, _ := uuid.NewV4()
    prd.Id = uuid.String()
    

    fmt.Println( "Product:", prd )
    key := prd.Id
    value, _ := json.Marshal(prd)
    err = client.Set(key, value, 0).Err()
    if err != nil {
      fmt.Println(err)
      formatter.JSON(w, http.StatusInternalServerError, err)
      return
    }

    formatter.JSON(w, http.StatusOK, prd)
  }
}

// API Read Order(s)
func starbucksProductDetailsHandler(formatter *render.Render) http.HandlerFunc {
  return func(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    var uuid string = params["id"]
    if uuid == "" {
      keys := client.Keys("*")
      if keys == nil {
        fmt.Println("Product not found")
        formatter.JSON(w, http.StatusOK, nil)
        return
      }
      var products []product
      for _, value := range keys.Val() {
        var p product
        var prd, err = client.Get(value).Result()
        if err != nil {
          continue
        }
        fmt.Println("Product details:", prd)
        json.Unmarshal([]byte(prd), &p)
        products = append(products, p)
      }
      formatter.JSON(w, http.StatusOK, products)
    } else {
      var p product
      var prd, err = client.Get(uuid).Result()
      if err != nil {
        fmt.Println("Order not found")
        formatter.JSON(w, http.StatusOK, nil)
        return
      }
      fmt.Println("Product details:", prd)
      json.Unmarshal([]byte(prd), &p)
      formatter.JSON(w, http.StatusOK, prd)
    }
  }
}

// API Delete products
func starbucksDeleteProductHandler(formatter *render.Render) http.HandlerFunc {
  return func(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    var uuid string = params["id"]
    var _, err = client.Get(uuid).Result()
    if err != nil {
      fmt.Println("Product not found")
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
