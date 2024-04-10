package input

import (
	"bufio"
	"fmt"
	"log/slog"
	"mintalk/server/db"
	"os"
)

type Console struct {
	database *db.Connection
}

func NewConsole(database *db.Connection) *Console {
	return &Console{database}
}

func (console *Console) InputLoop() {
	for {
		if err := console.Input(); err != nil {
			slog.Error("error executing command", "err", err)
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
	return console.Execute(line)
}

func (console *Console) Execute(rawCommand string) error {
	command := console.ParseCommand(rawCommand)
	if len(command) == 0 {
		return nil
	}
	switch command[0] {
	case "op":
		return console.op(command[1:])
	case "deop":
		return console.deop(command[1:])
	case "useradd":
		return console.useradd(command[1:])
	case "userdel":
		return console.userdel(command[1:])
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
