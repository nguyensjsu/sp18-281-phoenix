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
  mx.HandleFunc("/product", starbucksNewProductHandler(formatter)).Methods("POST")
  mx.HandleFunc("/product/{id}", starbucksProductDetailsHandler(formatter)).Methods("GET")
  mx.HandleFunc("/product/{id}", starbucksProductUpdateHandler(formatter)).Methods("PUT")
}

// Starbucks API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
  return func(w http.ResponseWriter, req *http.Request) {
    formatter.JSON(w, http.StatusOK, struct{ Ping string }{"Pong"})
  }
}

// API Create New starbucks Product
func starbucksNewProductHandler(formatter *render.Render) http.HandlerFunc {
  return func(w http.ResponseWriter, req *http.Request) {
    var prd order
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

// API Get Product Details 
func starbucksProductDetailsHandler(formatter *render.Render) http.HandlerFunc {
  return func(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    var uuid string = params["id"]
    if uuid == "" {
      keys := client.Keys("*")
      if keys == nil {
        fmt.Println("Product not found.")
        formatter.JSON(w, http.StatusNotFound, nil)
        return
      }
      var prod_array []string
      for key, value := range keys.Val() {
        fmt.Println("Key:", key, "Product:", value)
        prod_array = append(prod_array, value)
      }
      formatter.JSON(w, http.StatusOK, prod_array)
    } else {
      var prd, err = client.Get(uuid).Result()
      if err != nil {
        fmt.Println("product not found.")
        formatter.JSON(w, http.StatusNotFound, err)
        return
      }
      fmt.Println("Product: ", prd)
      formatter.JSON(w, http.StatusOK, prd)
    }
   }
 }



//API for update Product
func starbucksProductUpdateHandler (formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
      params := mux.Vars(req)
      var uuid string = params["id"]
      var prd order
      err := json.NewDecoder(req.Body).Decode(&prd)
      if err != nil {
                      fmt.Println(err)
                      formatter.JSON(w, http.StatusInternalServerError, err)
                      return
                    }
      if uuid == "" {
                      keys := client.Keys("*")
                      if keys == nil {
                      fmt.Println("Product not found.")
                      formatter.JSON(w, http.StatusNotFound, nil)
                      return
                                      }
                    } else {
                              var prd, err = client.Get(uuid).Result()
                              if err != nil {
                              fmt.Println("Product not found.")
                              formatter.JSON(w, http.StatusNotFound, err)
                              return
                                            }
                      fmt.Println("initial Details: ", prd)
                      //formatter.JSON(w, http.StatusOK, ord)
                            }
        prd.Id = uuid
        
        key := prd.Id
        value, _ := json.Marshal(prd)
        err = client.Set(key, value, 0).Err()
        fmt.Println("updated Details: ", prd)
        formatter.JSON(w, http.StatusOK, prd)
  }
}



