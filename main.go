package main

import (
  "os"
  "github.com/gorilla/mux"
  "html/template"
)

func newRouter() *mux.Router {
  r := mux.NewRouter()
  
  r.HandleFunc("/", indexHandler)
  //STATIC FILES
  staticFileDirectory := http.Dir("./assets/")
  staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
  r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

  return r
}

func main() {
  router := newRouter()
  port := ":80"
  http.ListenAndServer(port, router)
}