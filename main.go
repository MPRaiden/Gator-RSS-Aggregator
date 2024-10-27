package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/MPRaiden/gator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	commandNames map[string]func(*state, command) error
}

func main() {
	// Read from config file and create a state struct that holds a pointer to the config file
	gatorConfig, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	s := state{
		config: &gatorConfig,
	}

	// Create a commands struct and register a login handler function on it
	cmds := commands{commandNames: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)

	// Get cmd line arguments
	if len(os.Args) < 2 {
		fmt.Println("Error: not enough arguments")
		os.Exit(1)
	}

	argsWithoutProgram := os.Args[1:]
	commandName := argsWithoutProgram[0]
	commandArguments := argsWithoutProgram[1:]

	// Create a command from the cmd line arguments & runs it
	cmd := command{
		name:      commandName,
		arguments: commandArguments,
	}
	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		os.Exit(1)
	}

	// Read file again after update
	gatorConfig, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandNames[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.commandNames[cmd.name]
	if !ok {
		return fmt.Errorf("func run(): provided command not registered in commands map!")
	}
	return handler(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("login command requires a username")
	}

	if err := s.config.SetUser(cmd.arguments[0]); err != nil {
		return fmt.Errorf("failed to set user: %w: ", err)
	}
	fmt.Println("User name: ", cmd.arguments[0], " has been set")

	return nil
}
