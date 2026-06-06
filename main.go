package main

import (
	"log"
	"meal-planner-service/pkg/db"
)

func main() {
	didTestSucceed := db.NewDBTest()
	if !didTestSucceed {
		log.Fatal("Database test failed")
	} else {
		log.Println("Database test succeeded")
	}
}
