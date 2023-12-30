package main

import (
	"log"

	"github.com/Ibraheem0x45/dbfz-rollback/server"
	"github.com/Ibraheem0x45/dbfz-rollback/utils"
)

func main() {
	go utils.RunDBFZ()

	rollback_proxy := server.New()
	log.Println("Server starting at 8080")

	if err := rollback_proxy.Run(":8080"); err != nil {
		log.Panicf("Failed to run server reason: %s\n", err)
	}
}
