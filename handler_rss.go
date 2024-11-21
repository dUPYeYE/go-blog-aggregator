package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/dUPYeYE/go-blog-aggregator/internal/database"
	"github.com/dUPYeYE/go-blog-aggregator/rss"
)

func scrapeFeeds(s *State) {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Fatalf("Error getting next feed to fetch: %v", err)
	}

	if err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID); err != nil {
		log.Fatalf("Error marking feed fetched: %v", err)
	}

	rssFeed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		log.Fatalf("Error fetching feed: %v", err)
	}

	printFeed(*rssFeed)

	for _, item := range rssFeed.Channel.Item {
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Fatalf("Error parsing time: %v", err)
			continue
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			FeedID:      nextFeed.ID,
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: publishedAt,
		})
		if err != nil {
			log.Fatalf("Error creating post: %v", err)
			continue
		}
	}
}

func printFeed(rssFeed rss.RSSFeed) {
	fmt.Printf("Title: %s\n", rssFeed.Channel.Title)

	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("Item title: %s\n", item.Title)
	}
}

func handlerAggregate(s *State, cmd Command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: agg <time_between_requests>")
	}

	duration, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error parsing duration: %v", err)
	}
	ticker := time.NewTicker(duration)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerNewFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Usage: addfeed <feed_title> <feed_url>")
	}

	feedTitle := cmd.args[0]
	feedURL := cmd.args[1]

	createArgs := database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      feedTitle,
		Url:       feedURL,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	feed, err := s.db.CreateFeed(context.Background(), createArgs)
	if err != nil {
		return fmt.Errorf("Error creating feed: %v", err)
	}
	fmt.Println("Added feed", feedTitle)

	createFollowArgs := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if _, err = s.db.CreateFeedFollow(context.Background(), createFollowArgs); err != nil {
		return fmt.Errorf("Error following feed: %v", err)
	}

	return nil
}

func handlerListFeeds(s *State, cmd Command, user database.User) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("Usage: feeds")
	}

	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting feeds: %v", err)
	}

	for id, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feeds[id].UserID)
		if err != nil {
			return fmt.Errorf("Error getting user: %v", err)
		}

		fmt.Printf("%s\n%s\n%s\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}

func handlerFollowFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: follow <feed_url>")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error getting feed: %v", err)
	}

	createArgs := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), createArgs)
	if err != nil {
		return fmt.Errorf("Error following feed: %v", err)
	}

	fmt.Printf("%s, you are now following \"%s\".\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerUnfollowFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: unfollow <feed_url>")
	}

	if err := s.db.RemoveFeedFollow(context.Background(), database.RemoveFeedFollowParams{
		Name: user.Name,
		Url:  cmd.args[0],
	}); err != nil {
		return fmt.Errorf("Error unfollowing feed: %v", err)
	}

	fmt.Printf("%s, you are no longer following \"%s\".\n", user.Name, cmd.args[0])
	return nil
}

func handlerGetFollowsForUser(s *State, cmd Command, user database.User) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("Usage: following")
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Error getting follows: %v", err)
	}

	for _, follow := range follows {
		fmt.Printf("%s\n", follow.FeedName)
	}

	return nil
}

func handlerBrowsePosts(s *State, cmd Command, user database.User) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("Usage: browse <limit>(optional)")
	}

	var limit int32 = 3
	if len(cmd.args) == 1 {
		l, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("Error parsing limit: %v", err)
		}
		limit = int32(l)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("Error getting posts: %v", err)
	}

	for _, post := range posts {
		fmt.Printf("%s\n%s\n%s\n", post.Title, post.Url, post.Description)
	}

	return nil
}
