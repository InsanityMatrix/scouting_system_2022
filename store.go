package main

import (
	"database/sql"
	"fmt"
	"strconv"
)

type dbStore struct {
	db *sql.DB
}

func (store *dbStore) maintainDatabase() {
	store.db.Query("DELETE FROM scouting a USING scouting b WHERE a.match = b.match AND a.team = b.team")
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
