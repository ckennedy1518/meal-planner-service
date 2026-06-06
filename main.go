package main

import (
	"log"
	"meal-planner-service/pkg/db"
)

func main() {
	didTestSucceed := db.NewDirectDBTest()
	if didTestSucceed {
		log.Println("Database test succeeded")
	} else {
		log.Fatal("Database test failed")
	}
}
