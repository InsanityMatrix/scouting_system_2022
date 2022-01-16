package main

import (
  "github.com/gorilla/mux"
  "net/http"
  "html/template"
)

func newRouter() *mux.Router {
  r := mux.NewRouter()
  
  r.HandleFunc("/", indexHandler)
  r.HandleFunc("/test", gameTestHandler)
  //STATIC FILES
  staticFileDirectory := http.Dir("./assets/")
  staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
  r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

  return r
}

func main() {
  router := newRouter()
  port := ":80"
  http.ListenAndServe(port, router)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html")
  tmpl, err := template.ParseFiles("index.html")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  tmpl.Execute(w, map[string]string{})
}

func gameTestHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html")
  tmpl, err := template.ParseFiles("gameTest.html")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  tmpl.Execute(w, map[string]string{})
}