package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/MPRaiden/gator/internal/config"
	"github.com/MPRaiden/gator/internal/database"

	_ "github.com/lib/pq"
)

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

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
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
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("feeds", handlerListFeeds)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

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

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Error while retrieving next feed to fetch from database %v", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return fmt.Errorf("Error while marking feed fetched %v", err)
	}

	rssFeed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("Error while retrieving feed by url %v", err)
	}

	for _, item := range rssFeed.Channel.Items {
		postArgs := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{String: item.Title, Valid: item.Title != ""},
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			FeedID:      uuid.NullUUID{UUID: nextFeed.ID, Valid: true},
		}

		// Parse PubDate into time.Time
		parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate) // Adjust format if needed
		if err != nil {
			return fmt.Errorf("Error parsing time %v", err)
		}
		postArgs.PublishedAt = parsedTime

		err = s.db.CreatePost(context.Background(), postArgs)
		if err != nil {
			return fmt.Errorf("Error while creating database post %v", err)
		}
	}

	return nil
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	rssFeed := &RSSFeed{}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request with context: %v", err)
	}

	req.Header.Add("User-Agent", "gator")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %v", err)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading from request body: %v", err)
	}
	res.Body.Close()

	err = xml.Unmarshal(data, rssFeed)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling data: %v", err)
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	for i := range rssFeed.Channel.Items {
		rssFeed.Channel.Items[i].Title = html.UnescapeString(rssFeed.Channel.Items[i].Title)
		rssFeed.Channel.Items[i].Description = html.UnescapeString(rssFeed.Channel.Items[i].Description)
	}

	return rssFeed, nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("Less than 1 command argument provided")
	}
	url := cmd.arguments[0]

	userID := uuid.NullUUID{
		UUID:  user.ID,
		Valid: true,
	}
	delFeedFollowArgs := database.DeleteFeedFollowParams{
		UserID: userID,
		Url:    url,
	}
	err := s.db.DeleteFeedFollow(context.Background(), delFeedFollowArgs)
	if err != nil {
		return fmt.Errorf("Error while deleting feed follow %v", err)
	}

	return nil
}

func handlerFollowing(s *state, _ command, user database.User) error {
	userID := uuid.NullUUID{
		UUID:  user.ID,
		Valid: true,
	}
	userFeedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), userID)
	if err != nil {
		return fmt.Errorf("couldn't get feed follows: %v", err)
	}

	for _, userFeed := range userFeedFollows {
		fmt.Printf("* %v\n", userFeed.FeedName)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("Less than 1 command argument provided")
	}
	url := cmd.arguments[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	userID := uuid.NullUUID{
		UUID:  user.ID,
		Valid: true,
	}
	feedID := uuid.NullUUID{
		UUID:  feed.ID,
		Valid: true,
	}
	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    userID,
		FeedID:    feedID,
	}

	feedFollows, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("Error while creating feed follow %v", err)
	}

	if len(feedFollows) > 1 {
		feedFollow := feedFollows[0]
		fmt.Printf("Feed '%v' followed by user '%v'\n", feedFollow.FeedName, feedFollow.UserName)
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("agg command requires exactly one argument: time_between_reqs")
	}
	timeFromCmd := cmd.arguments[0]

	timeBetweenRequests, err := time.ParseDuration(timeFromCmd)
	if err != nil {
		return fmt.Errorf("Error while parsing time duration %v", err)
	}

	fmt.Printf("Collecting feed every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return fmt.Errorf("Error while scrapping feeds %v", err)
		}
	}
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

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("Less than 2 arguments provided!")
	}
	feedName := cmd.arguments[0]
	feedURL := cmd.arguments[1]

	userID := uuid.NullUUID{
		UUID:  user.ID,
		Valid: true,
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    userID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("Error while creating new feed: %v", err)
	}

	addFeed(feed.Name, feed.Url)

	feedID := uuid.NullUUID{
		UUID:  feed.ID,
		Valid: true,
	}
	feedFollowArgs := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
		FeedID:    feedID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), feedFollowArgs)
	if err != nil {
		return fmt.Errorf("Error while creating feed follow %v", err)
	}

	return nil
}

func addFeed(name, url string) {
	fmt.Printf("Creating new feed %s with the following url %s", name, url)
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error while getting users from database %v", err)
	}

	for _, user := range users {
		if s.cfg.CurrentUserName == user {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error while retrieving feeds from database %v", err)
	}
	for _, feed := range feeds {
		fmt.Printf("* %v\n", feed.Name)
		fmt.Printf("* %v\n", feed.Url)
		fmt.Printf("* %v\n", feed.Name_2)
	}
	return nil
}
