package main

import (
	"context"
	"fmt"
)

func handlerReset(s *State, cmd Command) error {
	if s.db.RemoveAllUsers(context.Background()) != nil {
		return fmt.Errorf("Error deleting all users")
	}
	fmt.Println("Removed all users successfully")

	if s.db.RemoveAllFeeds(context.Background()) != nil {
		return fmt.Errorf("Error deleting all feeds")
	}
	fmt.Println("Removed all feeds successfully")

	return nil
}
