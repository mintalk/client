package input

import (
	"bufio"
	"fmt"
	"log/slog"
	"mintalk/server/db"
	"mintalk/server/network"
	"os"
	"strings"
)

type Console struct {
	database *db.Connection
	server   *network.Server
}

func NewConsole(database *db.Connection, server *network.Server) *Console {
	return &Console{database, server}
}

func (console *Console) InputLoop() {
	for {
		if err := console.Input(); err != nil {
			slog.Warn("error executing command", "err", err)
		}
	}
}

func (console *Console) Input() error {
	fmt.Print("> ")
	in := bufio.NewReader(os.Stdin)
	line, err := in.ReadString('\n')
	if err != nil {
		return err
	}
	line = strings.TrimSuffix(line, "\n")
	return console.Execute(line)
}

func (console *Console) Execute(rawCommand string) error {
	command := console.ParseCommand(rawCommand)
	if len(command) == 0 {
		return nil
	}
	switch command[0] {
	case "user":
		return console.user(command[1:])
	case "group":
		return console.group(command[1:])
	case "channel":
		return console.channel(command[1:])
	default:
		return fmt.Errorf("command not found: %s", command[0])
	}
}

func (console *Console) ParseCommand(rawCommand string) []string {
	command := make([]string, 0)
	inQuote := false
	currentPart := ""
	for i := 0; i < len(rawCommand); i++ {
		char := rawCommand[i]
		if char == ' ' && !inQuote {
			command = append(command, currentPart)
			currentPart = ""
			continue
		}
		if char == '"' {
			inQuote = !inQuote
			continue
		}
		currentPart += string(char)
	}
	if currentPart != "" {
		command = append(command, currentPart)
	}
	return command
}
