package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func RunDBFZ() {
	entries, err := os.ReadDir("./")

	if err != nil {
		log.Printf("Failed to read current directory reason: %s\n", err)
		return
	}

	game_file := ""

	for _, entry := range entries {
		if entry.Type().IsRegular() && strings.Contains(entry.Name(), "-eac-nop-loaded.exe") {
			game_file = entry.Name()
			break
		}
	}

	if len(game_file) <= 0 {
		log.Println("Failed to find game executable! Make sure game excutable contains '-eac-nop-loaded'")
		return
	}

	command := exec.Command(fmt.Sprintf(`.\%s`, game_file))

	log.Printf("Launching %s\n", game_file)

	if err = command.Run(); err != nil {
		log.Printf("Failed to launch %s reason: %s\n", game_file, err)
	}
}
