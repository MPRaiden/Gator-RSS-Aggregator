package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"time"

	"github.com/MPRaiden/gator/internal/config"
	"github.com/MPRaiden/gator/internal/database"
)
import _ "github.com/lib/pq"

type state struct {
	cfg *config.Config
	db  *database.Queries
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

	db, err := sql.Open("postgres", gatorConfig.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	queries := database.New(db)

	s := state{
		cfg: &gatorConfig,
		db:  queries,
	}

	// Create a commands struct and register a login handler function on it
	cmds := commands{commandNames: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerResetDB)

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

	loginName := cmd.arguments[0]
	_, err := s.db.GetUser(context.Background(), loginName)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("User does not exist in database")
	} else if err != nil {
		return fmt.Errorf("failed to query user: %w", err)
	}

	if err := s.cfg.SetUser(cmd.arguments[0]); err != nil {
		return fmt.Errorf("failed to set user: %w: ", err)
	}
	fmt.Println("User name: ", cmd.arguments[0], " has been set")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return errors.New("register command requires a name")
	}

	name := cmd.arguments[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err == nil {
		return fmt.Errorf("user already exists in database")
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to query user: %w", err)
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	newUser, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("failed to create new user: %w", err)
	}

	s.cfg.SetUser(name)

	fmt.Printf("Successfully created new user: %s\n", name)
	log.Printf("User details: %+v\n", newUser)

	return nil
}

func handlerResetDB(s *state, cmd command) error {
	err := s.db.DelUsers(context.Background())
	if err != nil {
		return errors.New("error while reseting database")
	}
	log.Printf("Database reset successfull!")
	return nil
}
