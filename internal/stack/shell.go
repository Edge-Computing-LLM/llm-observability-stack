package stack

import (
	"strings"
)

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}

func joinCommands(commands []string) string {
	var filtered []string
	for _, command := range commands {
		if strings.TrimSpace(command) != "" {
			filtered = append(filtered, command)
		}
	}
	return strings.Join(filtered, "\n")
}
