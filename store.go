package main

import (
	"database/sql"
	"fmt"
	"strconv"
)

type dbStore struct {
	db *sql.DB
}

func mergeSort(fp []TeamPoints, sp []TeamPoints) []TeamPoints {
	var n = make([]TeamPoints, len(fp)+len(sp))

	var fpIndex = 0
	var spIndex = 0

	var nIndex = 0

	for fpIndex < len(fp) && spIndex < len(sp) {
		if fp[fpIndex].Points < sp[spIndex].Points {
			n[nIndex] = fp[fpIndex]
			fpIndex++
		} else {
			n[nIndex] = sp[spIndex]
			spIndex++
		}

		nIndex++
	}

	for fpIndex < len(fp) {
		n[nIndex] = fp[fpIndex]
		fpIndex++
		nIndex++
	}
	for spIndex < len(sp) {
		n[nIndex] = sp[spIndex]
		spIndex++
		nIndex++
	}
	return n
}
func sortPointList(arr []TeamPoints) []TeamPoints {
	if len(arr) == 1 {
		return arr
	}
	var fp = sortPointList(arr[0 : len(arr)/2])
	var sp = sortPointList(arr[len(arr)/2:])

	return mergeSort(fp, sp)
}
func (store *dbStore) getAllTeams() []int {
	var allTeams []int
	rows, err := store.db.Query("SELECT team FROM scouting;")
	if err != nil {
		fmt.Println(err)
		return allTeams
	}
	defer rows.Close()

	for rows.Next() {
		var team int
		err = rows.Scan(&team)
		if err != nil {
			fmt.Println(err)
			return allTeams
		}
		exists := false
		for _, t := range allTeams {
			if t == team {
				exists = true
			}
		}
		if !exists {
			allTeams = append(allTeams, team)
		}
	}

	return allTeams
}

func (store *dbStore) getBestAuton(teams []int) []TeamPoints {
	//Make List to store teams and their average auton points
	t := []TeamPoints{}
	for _, team := range teams {
		data, auton, _ := store.getTeamData(team)
		//Go through auton pts array and add up all points
		points := 0
		for _, shot := range auton {
			if shot.Result == "topbasket" {
				points += 4
			} else if shot.Result == "bottombasket" {
				points += 2
			}
		}
		//Get Average Points per round
		avgPoints := float64(points) / float64(len(data))
		nTeam := TeamPoints{Team: team, Points: avgPoints}
		t = append(t, nTeam)
	}
	//Sort so that most auton points are at top
	t = sortPointList(t)
	return t
}
func (store *dbStore) getTeamData(team int) ([]TeamData, []Shot, []Shot) {
	rows, err := store.db.Query("SELECT match,alliancestation,preloaded,movedstart,topintake,floorintake,attemptedlower,attemptedmiddle,attemptedhigh,attemptedtraversal,successful,endgamecomment,defense,attempted,disconnected,comments FROM scouting WHERE team=$1", team)
	if err != nil {
		fmt.Println(err)
		return []TeamData{}, []Shot{}, []Shot{}
	}
	defer rows.Close()

	teamData := []TeamData{}
	for rows.Next() {
		data := TeamData{}

		err = rows.Scan(&data.Match, &data.AllianceStation, &data.Preloaded, &data.MovedStart, &data.TopIntake, &data.FloorIntake, &data.AttemptedLower, &data.AttemptedMiddle, &data.AttemptedHigh, &data.AttemptedTraversal, &data.Successful, &data.EndgameComment, &data.Defense, &data.Attempted, &data.Disconnected, &data.Comments)
		data.Team = team
		if err != nil {
			fmt.Println(err)
			return []TeamData{}, []Shot{}, []Shot{}
		}
		teamData = append(teamData, data)
	}

	autonShots := []Shot{}
	teleopShots := []Shot{}
	for _, entry := range teamData {
		//Get Auton Shots
		rows, err = store.db.Query("SELECT X,Y,Result FROM auton_" + strconv.Itoa(entry.Match) + "_" + strconv.Itoa(entry.Team) + ";")
		if err != nil {
			fmt.Println(err)
			return teamData, []Shot{}, []Shot{}
		}
		defer rows.Close()

		for rows.Next() {
			shot := Shot{}
			err = rows.Scan(&shot.X, &shot.Y, &shot.Result)
			if err != nil {
				fmt.Println(err)
			}
			autonShots = append(autonShots, shot)
		}
		//Get Teleop Shots
		rows, err = store.db.Query("SELECT X,Y,Result FROM teleop_" + strconv.Itoa(entry.Match) + "_" + strconv.Itoa(entry.Team) + ";")
		if err != nil {
			fmt.Println(err)
			return teamData, []Shot{}, []Shot{}
		}
		defer rows.Close()

		for rows.Next() {
			shot := Shot{}
			err = rows.Scan(&shot.X, &shot.Y, &shot.Result)
			if err != nil {
				fmt.Println(err)
			}
			teleopShots = append(teleopShots, shot)
		}
	}

	return teamData, autonShots, teleopShots
}
func (store *dbStore) logScout(match int, team int, allianceStation string, preloaded bool,
	movedStart bool, topIntake bool, floorIntake bool, attemptedLower bool,
	attemptedMiddle bool, attemptedHigh bool, attemptedTraversal bool, successful int,
	endgameComment string, defense bool, attempted bool, disconnected bool, comments string, autonShots []Shot, teleopShots []Shot) {

	//Insert into database
	//returns rows, err
	_, err := store.db.Query("INSERT INTO scouting(match, team, alliancestation, preloaded, movedstart, topintake, floorintake,attemptedlower,attemptedmiddle, attemptedhigh,attemptedtraversal,successful,endgamecomment,defense,attempted,disconnected,comments) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)",
		match, team, allianceStation, preloaded, movedStart, topIntake, floorIntake, attemptedLower,
		attemptedMiddle, attemptedHigh, attemptedTraversal, successful, endgameComment, defense, attempted, disconnected, comments)
	if err != nil {
		fmt.Println(err)
		return
	}

	autonShotsTB := "auton_" + strconv.Itoa(match) + "_" + strconv.Itoa(team)
	var exists bool
	row := store.db.QueryRow("SELECT EXISTS ( SELECT FROM information_schema.tables WHERE table_schema='public' AND table_name=$1);", autonShotsTB)
	err = row.Scan(&exists)

	if err != nil {
		fmt.Println(err)
	}

	if exists {
		store.db.Query("DELETE FROM " + autonShotsTB + " WHERE 1=1;")
	} else {
		store.db.Query("CREATE TABLE " + autonShotsTB + " (id SERIAL PRIMARY KEY, X decimal, Y decimal, Result text);")
	}
	//To insert list into table
	autonQuery := "INSERT INTO " + autonShotsTB + " (X,Y,Result) VALUES "
	vals := []interface{}{}
	ticker := 1
	for _, shot := range autonShots {
		autonQuery += fmt.Sprintf("($%d, $%d, $%d),", ticker, ticker+1, ticker+2)
		ticker += 3
		vals = append(vals, shot.X, shot.Y, shot.Result)
	}
	autonQuery = autonQuery[0 : len(autonQuery)-1]
	fmt.Println(autonQuery)
	stmt, err := store.db.Prepare(autonQuery)
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec(vals...)
	if err != nil {
		fmt.Println(err)
	}

	//TELEOP
	teleopShotsTB := "teleop_" + strconv.Itoa(match) + "_" + strconv.Itoa(team)
	row = store.db.QueryRow("SELECT EXISTS ( SELECT FROM information_schema.tables WHERE table_schema='public' AND table_name=$1);", teleopShotsTB)
	err = row.Scan(&exists)

	if err != nil {
		fmt.Println(err)
	}

	if exists {
		store.db.Query("DELETE FROM " + teleopShotsTB + " WHERE 1=1;")
	} else {
		store.db.Query("CREATE TABLE " + teleopShotsTB + " (id SERIAL PRIMARY KEY, X decimal, Y decimal, Result text);")
	}
	//To insert list into table
	teleopQuery := "INSERT INTO " + teleopShotsTB + " (X, Y, Result) VALUES "
	vals = []interface{}{}
	ticker = 1
	for _, shot := range teleopShots {
		teleopQuery += fmt.Sprintf("($%d, $%d, $%d),", ticker, ticker+1, ticker+2)
		ticker += 3
		vals = append(vals, shot.X, shot.Y, shot.Result)
	}
	teleopQuery = teleopQuery[0 : len(teleopQuery)-1]
	stmt, _ = store.db.Prepare(teleopQuery)
	stmt.Exec(vals...)
}

//ESSENTIALS:
var store dbStore

func InitStore(s dbStore) {
	store = s
}
