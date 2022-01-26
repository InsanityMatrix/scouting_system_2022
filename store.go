package main

import (
	"database/sql"
	"strconv"
	"fmt"
)

type dbStore struct {
	db *sql.DB
}


func (store *dbStore) logScout(match int, team int, allianceStation string, preloaded bool,
	movedStart bool, topIntake bool, floorIntake bool, attemptedLower bool,
	attemptedMiddle bool, attemptedHigh bool, attemptedTraversal bool, successful int,
	endgameComment string, defense bool, attempted bool, disconnected bool, comments string, autonShots []Shot, teleopShots []Shot) {
	
	//Insert into database
	//returns rows, err
	_, err := store.db.Query("INSERT INTO scouting(match, team, alliancestation, preloaded, movedstart, topintake, floorintake,attemptedlower,attemptedmiddle, attemptedhigh,attemptedtraversal,successful,endgamecomment,defense,attempted,disconnected,comments) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)",
	match,team,allianceStation,preloaded,movedStart,topIntake,floorIntake,attemptedLower,
	attemptedMiddle,attemptedHigh,attemptedTraversal,successful,endgameComment,defense,attempted,disconnected,comments)
	if err != nil {
		fmt.Println(err)
	}

	autonShotsTB := "auton_" + strconv.Itoa(match) + "_" + strconv.Itoa(team)
	var exists bool
	row := store.db.QueryRow("SELECT EXISTS ( SELECT FROM information_schema.tables WHERE table_schema='public' AND table_name=$1);", autonShotsTB)
	err = row.Scan(&exists)

	if err != nil {
		fmt.Println(err)
	}
	
	if exists {
		store.db.Query("REMOVE FROM " + autonShotsTB + " WHERE 1=1;")
	} else {
		store.db.Query("CREATE TABLE " + autonShotsTB + " (id SERIAL PRIMARY KEY, X decimal, Y decimal, Result text);")
	}
	//To insert list into table
	for _, shot := range autonShots {
		store.db.Query("INSERT INTO " + autonShotsTB + "(X,Y,Result) VALUES ($1, $2, $3)", shot.X, shot.Y, shot.Result)
	}
	//TELEOP
	teleopShotsTB := "teleop_" + strconv.Itoa(match) + "_" + strconv.Itoa(team)
	row = store.db.QueryRow("SELECT EXISTS ( SELECT FROM information_schema.tables WHERE table_schema='public' AND table_name=$1);", teleopShotsTB)
	err = row.Scan(&exists)

	if err != nil {
		fmt.Println(err)
	}
	
	if exists {
		store.db.Query("REMOVE FROM " + teleopShotsTB + " WHERE 1=1;")
	} else {
		store.db.Query("CREATE TABLE " + teleopShotsTB + " (id SERIAL PRIMARY KEY, X decimal, Y decimal, Result text);")
	}
	//To insert list into table
	for _, shot := range teleopShots {
		store.db.Query("INSERT INTO " + teleopShotsTB + "(X,Y,Result) VALUES ($1, $2, $3)", shot.X, shot.Y, shot.Result)
	}
}


//ESSENTIALS:
var store dbStore
func InitStore(s dbStore) {
	store = s
}