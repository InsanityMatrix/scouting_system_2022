package main

import (
  "github.com/gorilla/mux"
  "net/http"
  "html/template"
  "database/sql"
  "strconv"
  "fmt"
  "time"
  _ "github.com/lib/pq"
)

type Shot struct {
  X float64 `json:"X"`
  Y float64 `json:"Y"`
  Result string `json:"Result"`
}

func newRouter() *mux.Router {
  r := mux.NewRouter()
  
  r.HandleFunc("/", indexHandler)
  r.HandleFunc("/test", gameTestHandler)
  r.HandleFunc("/submit", submitScoutHandler)
  //STATIC FILES
  staticFileDirectory := http.Dir("./assets/")
  staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
  r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

  return r
}

func main() {
  router := newRouter()
  port := ":80"

  url := "host=localhost port=5432 user=techhounds password=team868 dbname=scouting sslmode=disable"
  db, err := sql.Open("postgres", url)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    panic(err)
  }
  fmt.Println("Database open!")
  db.SetMaxOpenConns(25)
  db.SetMaxIdleConns(10)
  db.SetConnMaxLifetime(time.Hour)
  InitStore(dbStore{db:db})
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
func submitScoutHandler(w http.ResponseWriter, r *http.Request) {
  //Will parse form, log, redirect

  err := r.ParseForm()
  fmt.Println("Submit accessed")
  if err != nil {
    fmt.Println(fmt.Errorf("Error: %v", err))
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  //Get all Variables
  match, _ := strconv.Atoi(r.Form.Get("match"))
  team, _ := strconv.Atoi(r.Form.Get("team"))
  allianceStation := r.Form.Get("allianceStation")
  preloaded, _ := strconv.ParseBool(r.Form.Get("preloaded"))
  movedStart, _ := strconv.ParseBool(r.Form.Get("moveStart"))
  topIntake, _ := strconv.ParseBool(r.Form.Get("topIntake"))
  floorIntake, _ := strconv.ParseBool(r.Form.Get("floorIntake"))
  attemptedLower, _ := strconv.ParseBool(r.Form.Get("attemptedLower"))
  attemptedMiddle, _ := strconv.ParseBool(r.Form.Get("attemptedMiddle"))
  attemptedHigh, _ := strconv.ParseBool(r.Form.Get("attemptedHigh"))
  attemptedTraversal, _ := strconv.ParseBool(r.Form.Get("attemptedTraversal"))
  successful, _ := strconv.Atoi(r.Form.Get("successful"))
  endgameComment := r.Form.Get("endgameComment")
  defense, _ := strconv.ParseBool(r.Form.Get("defense"))
  attempted, _ := strconv.ParseBool(r.Form.Get("attempted"))
  disconnected, _ := strconv.ParseBool(r.Form.Get("disconnected"))
  comments := r.Form.Get("comments")

  aShotsL, _  := strconv.Atoi(r.Form.Get("ashotLength"))
  var autonShots []Shot
  for i := 0; i < aShotsL; i++ {
    var newShot Shot
    newShot.X, _ = strconv.ParseFloat(r.Form.Get("autonShots["+strconv.Itoa(i)+"][position][x]"),64)
    newShot.Y, _ = strconv.ParseFloat(r.Form.Get("autonShots["+strconv.Itoa(i)+"][position][y]"),64)
    newShot.Result = r.Form.Get("autonShots["+strconv.Itoa(i)+"][result]")
    fmt.Println(newShot)
    autonShots = append(autonShots, newShot)
  }
  fmt.Println(autonShots)
  tShotsL, _ := strconv.Atoi(r.Form.Get("tshotLength"))
  var teleopShots []Shot
  for i := 0; i < tShotsL; i++ {
    var newShot Shot
    newShot.X, _ = strconv.ParseFloat(r.Form.Get("teleopShots["+strconv.Itoa(i)+"][position][x]"),64)
    newShot.Y, _ = strconv.ParseFloat(r.Form.Get("teleopShots["+strconv.Itoa(i)+"][position][y]"),64)
    newShot.Result = r.Form.Get("teleopShots["+strconv.Itoa(i)+"][result]")
    fmt.Println(newShot)
    teleopShots = append(teleopShots, newShot)
  }

  fmt.Println(teleopShots)
  fmt.Println("Submitting to database")

  store.logScout(match, team, allianceStation, preloaded,movedStart,
  topIntake, floorIntake,attemptedLower,attemptedMiddle,attemptedHigh,
  attemptedTraversal,successful,endgameComment,defense,attempted,disconnected,comments, autonShots,teleopShots)
}