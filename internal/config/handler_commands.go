package config

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"

	"github.com/harsiz/goggregator/internal/database"
)

// Func: Login (1 arg) - logs in user via singular username
func HandlerLogin(s *State, cmd Command) error {
	var requiredArgs = 1
	if len(cmd.Args) < requiredArgs {
		log.Fatal("Too little arguments provided for given command.")
	} else if len(cmd.Args) > requiredArgs {
		log.Fatal("Too much arguments provided for given command.")
	}
	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}
	if err := s.Confg.SetUser(cmd.Args[0]); err != nil {
		return err
	}
	r, err := Read()
	if err != nil {
		return err
	}
	fmt.Printf("Username '%v' has been set.\n", r.CurrentUserName)
	return nil
}

// Func: Register (not to be confused with config_commands_overseer.go register)
// Registers given user
func HandlerRegister(s *State, cmd Command) error {
	var requiredArgs = 1
	if len(cmd.Args) < requiredArgs {
		log.Fatal("Too little arguments provided for given command.")
	} else if len(cmd.Args) > requiredArgs {
		log.Fatal("Too much arguments provided for given command.")
	}
	newUUID := uuid.New()
	userParams := database.CreateUserParams{
		ID: newUUID,
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Name: cmd.Args[0],
	}
	dbUser, err := s.Db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	if err = s.Confg.SetUser(dbUser.Name); err != nil {
		return err
	}
	fmt.Printf("User %s has been created!\n", dbUser.Name)
	return nil
}

// Func: Reset (no args) - deletes users from db
func HandlerReset(s *State, cmd Command) error {
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// Func: Users (no args) - Lists users and highlights current user
func HandlerUsers(s *State, cmd Command) error {
	users, err := s.Db.GetAllUsers(context.Background())
	if err != nil {
		return err
	}

	for _, u := range users {
		if s.Confg.CurrentUserName == u.Name {
			fmt.Printf("* %s (current)\n", u.Name)
		} else {
			fmt.Printf("* %s\n", u.Name)
		}
	}
	return nil
}

// Func: agg (1 arg) - Prints stuff from feed url
func HandlerAggregator(s *State, cmd Command) error {
	var requiredArgs = 1
	var feedURL string
	if len(cmd.Args) < requiredArgs {
		feedURL = "https://www.wagslane.dev/index.xml"
	} else if len(cmd.Args) > requiredArgs {
		log.Fatal("Too much arguments provided for given command.")
	}

	if feedURL != "https://www.wagslane.dev/index.xml" {
		feedURL = cmd.Args[0]
	}

	f, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", f)
	return nil
}

// Func: AddFeed (2 args - adds item to feed)
func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	var requiredArgs = 2
	if len(cmd.Args) < requiredArgs {
		log.Fatal("Too little arguments provided for given command.")
	} else if len(cmd.Args) > requiredArgs {
		log.Fatal("Too much arguments provided for given command.")
	}

	providedName := cmd.Args[0]
	providedURL := cmd.Args[1]

	feedParams := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		Name: sql.NullString{
			String: providedName,
			Valid: true,
		},
		Url: sql.NullString{
			String: providedURL,
			Valid: true,
		},
		UserID: user.ID,
	}

	feed, err := s.Db.CreateFeed(
		context.Background(),
		feedParams,
	)
	if err != nil {
		return err
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		UserID: user.ID,
		FeedID: feed.ID,
	}

	_, err = s.Db.CreateFeedFollow(
		context.Background(),
		feedFollowParams,
	)
	if err != nil {
		return err
	}
	fmt.Printf("New Feed Created: %v\nDetails: %+v\n", feed.Name.String, feed)
	return nil
}

// Func: Feeds (no args)
func HandlerFeeds(s *State, cmd Command) error {
	feeds, err := s.Db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, f := range feeds {
		user, err := s.Db.GetUserFromUserId(context.Background(), f.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Feed Name: %s\nFeed URL: %s\nUser Who Created: %s\n\n", f.Name.String, f.Url.String, user)
	}
	return nil
}

// Func: Follow (1 arg) - Creates a new feed follow record for user
func HandlerFollow(s *State, cmd Command, user database.User) error {
	var requiredArgs = 1
	if len(cmd.Args) < requiredArgs {
		log.Fatal("Too little arguments provided for given command. A URL is needed.")
	} else if len(cmd.Args) > requiredArgs {
		log.Fatal("Too much arguments provided for given command.")
	}

	
	feed, err := s.Db.GetFeedFromURL(
		context.Background(),
		sql.NullString{
			String: cmd.Args[0],
			Valid: true,
		},
	)
	if err != nil {
		return nil
	}
	feedFollowParams := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		UserID: user.ID,
		FeedID: feed.ID,
	}
	feedRow, err := s.Db.CreateFeedFollow(
		context.Background(),
		feedFollowParams,
	)
	if err != nil {
		return err
	}

	g, err := s.Db.GetFeedFromFeedId(context.Background(), feedRow.FeedID)
	if err != nil {
		return err
	}
	fmt.Printf("Feed: %s\nCurrent User: %s\n\n", g.Name.String, user.Name)
	return nil
}

// Func: Following (no args) - prints all feeds the user is following
func HandlerFollowing(s *State, cmd Command, user database.User) error {
	feedFollow, err := s.Db.GetFeedFollowsForUser(
		context.Background(),
		user.ID,
	)
	if err != nil {
		return err
	}

	for n, f := range feedFollow {
		feed, err := s.Db.GetFeedFromFeedId(context.Background(), f.FeedID)
		if err != nil {
			return err
		}
		fmt.Printf("Feed #%d: %s", n, feed.Name.String)
	}
	return nil
}

// Func: Unfollow (1 arg) - unfollows feedfollow
func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	var requiredArgs = 1
	if len(cmd.Args) < requiredArgs {
		log.Fatal("Too little arguments provided for given command. A URL is needed.")
	} else if len(cmd.Args) > requiredArgs {
		log.Fatal("Too much arguments provided for given command.")
	}

	url := cmd.Args[0]
	feed, err := s.Db.GetFeedFromURL(
		context.Background(),
		sql.NullString{
			String: url,
			Valid: true,
		},
	)
	if err != nil {
		return err
	}

	UnfollowParams := database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.Db.UnfollowFeed(
		context.Background(),
		UnfollowParams,
	)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully unfollowed Feed Pair: UID - %v FID - %v", user.ID, feed.ID)
	return nil
}