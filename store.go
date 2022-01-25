package main

import (
	"database/sql"
	"fmt"
)

type dbStore struct {
	db *sql.DB
}


func (store *dbStore) logScout(match int, team int, allianceStation string, preloaded bool,
	movedStart bool, topIntake bool, floorIntake bool, attemptedLower bool,
	attemptedMiddle bool, attemptedHigh bool, attemptedTraversal bool, successful int,
	endgameComment string, defense bool, attempted bool, disconnected bool, comments string) {
	
	//Insert into database
	//returns rows, err
	_, err := store.db.Query("INSERT INTO scouting(match, team, alliancestation, preloaded, movedstart, topintake, floorintake,attemptedlower,attemptedmiddle, attemptedhigh,attemptedtraversal,successful,endgamecomment,defense,attempted,disconnected,comments) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)",
	match,team,allianceStation,preloaded,movedStart,topIntake,floorIntake,attemptedLower,
	attemptedMiddle,attemptedHigh,attemptedTraversal,successful,endgameComment,defense,attempted,disconnected,comments)
	fmt.Println(err)
}


//ESSENTIALS:
var store dbStore
func InitStore(s dbStore) {
	store = s
}