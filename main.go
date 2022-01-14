package main

import (
  "os"
  "github.com/gorilla/mux"
  "net/http"
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html")
  tmpl, err := template.ParseFiles("templates/pregame.html")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  tmpl.Execute(w, map[string]string{})
}