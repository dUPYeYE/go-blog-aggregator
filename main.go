package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/dUPYeYE/go-blog-aggregator/internal/config"
	"github.com/dUPYeYE/go-blog-aggregator/internal/database"
)

type State struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	config, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		fmt.Println("Error opening database connection")
		return
	}
	defer db.Close()

	state := &State{
		db:  database.New(db),
		cfg: &config,
	}
	commands := Commands{
		cmds: make(map[string]func(*State, Command) error),
	}
	// For anyone
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	// For logged in users
	commands.register("users", middlewareLoggedIn(handlerListUsers))
	commands.register("agg", middlewareLoggedIn(handlerAggregate))
	commands.register("addfeed", middlewareLoggedIn(handlerNewFeed))
	commands.register("feeds", middlewareLoggedIn(handlerListFeeds))
	commands.register("follow", middlewareLoggedIn(handlerFollowFeed))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollowFeed))
	commands.register("following", middlewareLoggedIn(handlerGetFollowsForUser))
	commands.register("browse", middlewareLoggedIn(handlerBrowsePosts))

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: gator <command> [args...]")
		return
	}

	if err := commands.run(state, Command{name: args[1], args: args[2:]}); err != nil {
		log.Fatal(err)
	}
}
