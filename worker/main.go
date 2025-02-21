package main

import (
	"log"
	"task-queue/worker/queue"
	"task-queue/worker/db"
)

func main() {
	dbService, err := db.NewDbRequest()
	if err != nil {
		log.Println("error while creating the db service instance: ", err)
		return
	}
	dbConn, err := dbService.InitDB()
	if err != nil { 
		log.Println("error while creating the db connection: ", err)
		return
	}

	log.Println("Wroker started...")
	queue.Queue(dbConn)
}
