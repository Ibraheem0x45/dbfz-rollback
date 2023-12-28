package main

import (
	"log"

	"github.com/Ibraheem0x45/dbfz-rollback/server"
)

func main() {
	rollback_proxy := server.New()
	err := rollback_proxy.Run(":8080")

	if err != nil {
		log.Panicf("Failed to run server reason: %s\n", err)
		return
	}
}
