package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hop-/gotchat/internal/ui/tui/commands"
)

type ChatCommand = func(...string) tea.Cmd

var (
	chatCommands map[string]ChatCommand
)

func init() {
	// Initialize the chat commands map
	chatCommands = make(map[string]func(...string) tea.Cmd)

	// Add the exit command
	exitCommand := func(args ...string) tea.Cmd {
		return commands.Shutdown
	}
	chatCommands["exit"] = exitCommand
	chatCommands["quit"] = exitCommand
	chatCommands["q"] = exitCommand

	// Add the connect command
	connectCommand := func(args ...string) tea.Cmd {
		if len(args) < 1 || len(args) > 2 {
			return commands.Error("connect command requires host and port")
		}

		var host, port string

		if len(args) == 1 || strings.Contains(args[0], ":") {
			// If only one argument is provided, it should be in the format "host:port"
			parts := strings.SplitN(args[0], ":", 2)
			if len(parts) != 2 {
				return commands.Error("connect command requires host:port format")
			}
			host = parts[0]
			port = parts[1]
		} else if len(args) == 2 {
			host = args[0]
			port = args[1]
		}
		if host == "" || port == "" {
			return commands.Error("connect command requires non-empty host and port")
		}

		return commands.Connect(host, port)
	}
	chatCommands["connect"] = connectCommand
}

func chatCommandExecuted(name string, args ...string) tea.Cmd {
	// Check if the command exists
	if command, exists := chatCommands[name]; exists {
		return command(args...)
	}

	return commands.Error(fmt.Sprintf("chat command '%s' not found", name))
}

func AddChatCommand(name string, command ChatCommand) error {
	// Check a chat command by name
	if _, exists := chatCommands[name]; exists {
		return fmt.Errorf("chat command '%s' already exists", name)
	}

	chatCommands[name] = command

	return nil
}
