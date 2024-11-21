package main

import (
	"context"
	"fmt"

	"github.com/dUPYeYE/go-blog-aggregator/internal/database"
)

func middlewareLoggedIn(
	handler func(s *State, cmd Command, user database.User) error,
) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.Username)
		if err != nil {
			return fmt.Errorf("Error getting user: %v", err)
		}

		return handler(s, cmd, user)
	}
}
