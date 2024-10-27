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
		log.Fatal("Needs at least 2 command line arguments!")
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
		log.Fatal("error occured during the execution of cmds.run():", err)
	}

	// Read file again after update
	gatorConfig, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(gatorConfig.DBURL)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandNames[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	_, ok := c.commandNames[cmd.name]
	if !ok {
		return fmt.Errorf("func run(): provided command not registed in commands map!")
	}

	err := c.commandNames[cmd.name](s, cmd)
	if err != nil {
		return fmt.Errorf("func run(): error occured while running command: %q", cmd.name)
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("func handlerLogin(): command.arguments cannot be empty!")
	}

	if err := s.config.SetUser(cmd.arguments[0]); err != nil {
		return fmt.Errorf("failed to set user: %w: ", err)
	}
	fmt.Println("func handlerLogin(): user name: ", cmd.arguments[0], " has been set")

	return nil
}
