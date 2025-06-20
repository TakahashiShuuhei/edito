package command

import (
	"fmt"
	"sort"
	"strings"
)

type Handler func(args []string) error
type InteractiveHandler func(promptFunc func(prompt string) (string, error)) error

type Command struct {
	Name        string
	Description string
	Handler     Handler
	Interactive InteractiveHandler
	NeedsArgs   bool
}

type Registry struct {
	commands map[string]*Command
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]*Command),
	}
}

func (r *Registry) Register(name, description string, handler Handler) {
	r.commands[name] = &Command{
		Name:        name,
		Description: description,
		Handler:     handler,
	}
}

func (r *Registry) RegisterInteractive(name, description string, interactive InteractiveHandler) {
	r.commands[name] = &Command{
		Name:        name,
		Description: description,
		Interactive: interactive,
		NeedsArgs:   true,
	}
}

func (r *Registry) Execute(name string, args []string) error {
	cmd, exists := r.commands[name]
	if !exists {
		return fmt.Errorf("command not found: %s", name)
	}
	
	if cmd.Handler != nil {
		return cmd.Handler(args)
	}
	
	return fmt.Errorf("command %s needs interactive execution", name)
}

func (r *Registry) ExecuteInteractive(name string, promptFunc func(prompt string) (string, error)) error {
	cmd, exists := r.commands[name]
	if !exists {
		return fmt.Errorf("command not found: %s", name)
	}
	
	if cmd.Interactive != nil {
		return cmd.Interactive(promptFunc)
	}
	
	return fmt.Errorf("command %s is not interactive", name)
}

func (r *Registry) GetCommand(name string) *Command {
	return r.commands[name]
}

func (r *Registry) ListCommands() []*Command {
	commands := make([]*Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		commands = append(commands, cmd)
	}
	
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})
	
	return commands
}

func (r *Registry) SearchCommands(query string) []*Command {
	var results []*Command
	query = strings.ToLower(query)
	
	for _, cmd := range r.commands {
		if strings.Contains(strings.ToLower(cmd.Name), query) ||
		   strings.Contains(strings.ToLower(cmd.Description), query) {
			results = append(results, cmd)
		}
	}
	
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})
	
	return results
}