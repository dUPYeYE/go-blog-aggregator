package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/dUPYeYE/go-blog-aggregator/internal/database"
)

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: register <username>")
	}

	username := cmd.args[0]
	createArgs := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      username,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if exists, _ := s.db.GetUser(context.Background(), username); exists.Name == username {
		return fmt.Errorf("User %s already exists", username)
	}

	user, err := s.db.CreateUser(context.Background(), createArgs)
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}
	fmt.Println("Registered user", username)

	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("Error setting user: %w", err)
	}
	fmt.Printf("Logged in as %s:\n", username)
	fmt.Printf("ID: %s\n", user)

	return nil
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: login <username>")
	}

	if exists, _ := s.db.GetUser(context.Background(), cmd.args[0]); exists.Name != cmd.args[0] {
		return fmt.Errorf("User %s does not exist", cmd.args[0])
	}

	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("Error setting user: %w", err)
	}

	fmt.Println("Logged in as", cmd.args[0])
	return nil
}

func handlerListUsers(s *State, cmd Command, user database.User) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error listing users: %w", err)
	}
	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if user.Name == s.cfg.Username {
			fmt.Print(" (current)")
		}
		fmt.Printf("\n")
	}
	return nil
}
