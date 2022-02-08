package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Shot struct {
	X      float64 `json:"X"`
	Y      float64 `json:"Y"`
	Result string  `json:"Result"`
}
type TeamData struct {
	Match              int
	Team               int
	AllianceStation    string
	Preloaded          bool
	MovedStart         bool
	TopIntake          bool
	FloorIntake        bool
	AttemptedLower     bool
	AttemptedMiddle    bool
	AttemptedHigh      bool
	AttemptedTraversal bool
	Successful         int
	EndgameComment     string
	Defense            bool
	Attempted          bool
	Disconnected       bool
	Comments           string
}

type TeamOverview struct {
	Team        int
	Data        []TeamData
	AutonShots  []Shot
	TeleopShots []Shot
}

func newRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/data", dataHandler)
	r.HandleFunc("/submit", submitScoutHandler)
	r.HandleFunc("/team/{team}", teamDataHandler)
	r.HandleFunc("/overview", teamOverviewHandler)
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
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(30 * time.Second)
	InitStore(dbStore{db: db})

	http.ListenAndServe(port, router)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, map[string]string{})
}
func dataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tmpl, err := template.ParseFiles("data.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, map[string]string{})
}

type DatabaseResponse struct {
	Data   []TeamData
	Auton  []Shot
	Teleop []Shot
}

func teamDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	team, err := strconv.Atoi(vars["team"])
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	data, auton, teleop := store.getTeamData(team)
	responseData := DatabaseResponse{Data: data, Auton: auton, Teleop: teleop}

	jsonInfo, _ := json.Marshal(responseData)
	fmt.Fprint(w, string(jsonInfo))
	//Write JSON
}

type TeamStats struct {
	Team          int
	PercentTop    float64
	PercentBottom float64
	PercentMisses float64
	AmountTop     int
	AmountBottom  int
	AmountMisses  int
	AmountTrav    int
	AmountHigh    int
	AmountMiddle  int
	AmountLower   int
}

type ShotRanking struct {
	Team       int
	Percentage float64
	Total      int
}
type AmountRanking struct {
	Team   int
	Amount int
}
type TeamRankings struct {
	PercentTop    []ShotRanking
	PercentBottom []ShotRanking
	PercentMisses []ShotRanking
	AmountTrav    []AmountRanking
	AmountHigh    []AmountRanking
	AmountMiddle  []AmountRanking
	AmountLower   []AmountRanking
}

func teamOverviewHandler(w http.ResponseWriter, r *http.Request) {
	teams := store.getAllTeams()
	overviewList := []TeamOverview{}

	for _, team := range teams {
		//Compile Team Data List
		data, auton, teleop := store.getTeamData(team)
		o := TeamOverview{
			Team:        team,
			Data:        data,
			AutonShots:  auton,
			TeleopShots: teleop,
		}
		overviewList = append(overviewList, o)
	}
	//Now get stats of each team separately
	allStats := []TeamStats{}
	for _, team := range overviewList {
		totalShots := len(team.AutonShots) + len(team.TeleopShots)
		totalMisses := 0
		totalTop := 0
		totalBottom := 0

		for _, shot := range team.AutonShots {
			if shot.Result == "topbasket" {
				totalTop++
			} else if shot.Result == "bottombasket" {
				totalBottom++
			} else {
				totalMisses++
			}
		}
		for _, shot := range team.TeleopShots {
			if shot.Result == "topbasket" {
				totalTop++
			} else if shot.Result == "bottombasket" {
				totalBottom++
			} else {
				totalMisses++
			}
		}

		pTop := float64(totalTop) / float64(totalShots) * 100
		pBot := float64(totalBottom) / float64(totalShots) * 100
		pMis := float64(totalMisses) / float64(totalShots) * 100

		//Collect Hanging Totals
		trav := 0
		high := 0
		middle := 0
		lower := 0
		for _, entry := range team.Data {
			switch entry.Successful {
			case 1:
				lower++
			case 2:
				middle++
			case 3:
				high++
			case 4:
				trav++
			}
		}

		stats := TeamStats{
			Team:          team.Team,
			PercentTop:    pTop,
			PercentBottom: pBot,
			PercentMisses: pMis,
			AmountTop:     totalTop,
			AmountBottom:  totalBottom,
			AmountMisses:  totalMisses,
			AmountTrav:    trav,
			AmountHigh:    high,
			AmountMiddle:  middle,
			AmountLower:   lower,
		}
		allStats = append(allStats, stats)
	}

	//SORT TOP BASKET DATA
	rankingsTop := []ShotRanking{}
	for _, t := range allStats {
		rankingsTop = append(rankingsTop, ShotRanking{
			Team:       t.Team,
			Percentage: t.PercentTop,
			Total:      t.AmountTop,
		})
	}
	rankingsTop = sortShotList(rankingsTop, len(rankingsTop))
	rankingsBottom := []ShotRanking{}
	for _, t := range allStats {
		rankingsBottom = append(rankingsBottom, ShotRanking{
			Team:       t.Team,
			Percentage: t.PercentBottom,
			Total:      t.AmountBottom,
		})
	}
	rankingsBottom = sortShotList(rankingsBottom, len(rankingsBottom))
	rankingsMissed := []ShotRanking{}
	for _, t := range allStats {
		rankingsMissed = append(rankingsMissed, ShotRanking{
			Team:       t.Team,
			Percentage: t.PercentMisses,
			Total:      t.AmountMisses,
		})
	}
	rankingsMissed = sortShotList(rankingsMissed, len(rankingsMissed))

	rankingsTrav := []AmountRanking{}
	rankingsHigh := []AmountRanking{}
	rankingsMiddle := []AmountRanking{}
	rankingsLower := []AmountRanking{}
	for _, te := range allStats {
		t := AmountRanking{Team: te.Team, Amount: te.AmountTrav}
		h := AmountRanking{Team: te.Team, Amount: te.AmountHigh}
		m := AmountRanking{Team: te.Team, Amount: te.AmountMiddle}
		l := AmountRanking{Team: te.Team, Amount: te.AmountLower}

		rankingsTrav = append(rankingsTrav, t)
		rankingsHigh = append(rankingsHigh, h)
		rankingsMiddle = append(rankingsMiddle, m)
		rankingsLower = append(rankingsLower, l)
	}
	rankingsTrav = sortAmountList(rankingsTrav, len(rankingsTrav))
	rankingsHigh = sortAmountList(rankingsHigh, len(rankingsHigh))
	rankingsMiddle = sortAmountList(rankingsMiddle, len(rankingsMiddle))
	rankingsLower = sortAmountList(rankingsLower, len(rankingsLower))

	rankings := TeamRankings{
		PercentTop:    rankingsTop,
		PercentBottom: rankingsBottom,
		PercentMisses: rankingsMissed,
		AmountTrav:    rankingsTrav,
		AmountHigh:    rankingsHigh,
		AmountMiddle:  rankingsMiddle,
		AmountLower:   rankingsLower,
	}
	info, _ := json.Marshal(rankings)
	fmt.Fprint(w, string(info))
}
func sortAmountList(list []AmountRanking, n int) []AmountRanking {
	if n == 1 {
		return list
	}

	for i := 0; i < n-1; i++ {
		if list[i].Amount < list[i+1].Amount {
			temp := list[i]
			list[i] = list[i+1]
			list[i+1] = temp
		}
		list = sortAmountList(list, n-1)
	}
	return list
}
func sortShotList(list []ShotRanking, n int) []ShotRanking {
	if n == 1 {
		return list
	}

	for i := 0; i < n-1; i++ {
		if list[i].Percentage < list[i+1].Percentage {
			temp := list[i]
			list[i] = list[i+1]
			list[i+1] = temp
		}

		list = sortShotList(list, n-1)
	}
	return list
}
func submitScoutHandler(w http.ResponseWriter, r *http.Request) {
	//Will parse form, log, redirect

	err := r.ParseForm()
	fmt.Println("Submit accessed")
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
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

	if len(endgameComment) > 144 {
		endgameComment = endgameComment[0:143]
	}
	if len(comments) > 144 {
		comments = comments[0:143]
	}

	aShotsL, _ := strconv.Atoi(r.Form.Get("ashotLength"))
	var autonShots []Shot
	for i := 0; i < aShotsL; i++ {
		var newShot Shot
		newShot.X, _ = strconv.ParseFloat(r.Form.Get("autonShots["+strconv.Itoa(i)+"][position][x]"), 64)
		newShot.Y, _ = strconv.ParseFloat(r.Form.Get("autonShots["+strconv.Itoa(i)+"][position][y]"), 64)
		newShot.Result = r.Form.Get("autonShots[" + strconv.Itoa(i) + "][result]")
		autonShots = append(autonShots, newShot)
	}
	tShotsL, _ := strconv.Atoi(r.Form.Get("tshotLength"))
	var teleopShots []Shot
	for i := 0; i < tShotsL; i++ {
		var newShot Shot
		newShot.X, _ = strconv.ParseFloat(r.Form.Get("teleopShots["+strconv.Itoa(i)+"][position][x]"), 64)
		newShot.Y, _ = strconv.ParseFloat(r.Form.Get("teleopShots["+strconv.Itoa(i)+"][position][y]"), 64)
		newShot.Result = r.Form.Get("teleopShots[" + strconv.Itoa(i) + "][result]")
		teleopShots = append(teleopShots, newShot)
	}

	fmt.Println("Submitting to database")

	store.logScout(match, team, allianceStation, preloaded, movedStart,
		topIntake, floorIntake, attemptedLower, attemptedMiddle, attemptedHigh,
		attemptedTraversal, successful, endgameComment, defense, attempted, disconnected, comments, autonShots, teleopShots)
}
