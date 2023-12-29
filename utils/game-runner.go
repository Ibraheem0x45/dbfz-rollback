package utils

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

var (
	ErrGameExecutableNotFound = errors.New("game executable not found")
)

func RunDBFZ() (string, error) {
	entries, err := os.ReadDir("./")

	if err != nil {
		return "", err
	}

	game_file := ""

	for _, entry := range entries {
		if entry.Type().IsRegular() && strings.Contains(entry.Name(), "-eac-nop-loaded.exe") {
			game_file = entry.Name()
			break
		}
	}

	if len(game_file) <= 0 {
		return "", ErrGameExecutableNotFound
	}

	command := exec.Command("run", game_file)
	return game_file, command.Run()
}
