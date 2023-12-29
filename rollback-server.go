package main

import (
	"log"

	"github.com/Ibraheem0x45/dbfz-rollback/server"
	"github.com/Ibraheem0x45/dbfz-rollback/utils"
)

func main() {
	game_file, err := utils.RunDBFZ()

	if err != nil {
		log.Printf("Failed to run DBFZ rollback reason: %s\n", err)
	} else {
		log.Printf("Launching %s\n", game_file)
	}

	rollback_proxy := server.New()
	log.Println("Server started at 8080")

	err = rollback_proxy.Run(":8080")

	if err != nil {
		log.Panicf("Failed to run server reason: %s\n", err)
		return
	}
}
