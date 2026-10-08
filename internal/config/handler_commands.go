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

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}
	return nil
}

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

// 
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
func HandlerAddFeed(s *State, cmd Command) error {
	var requiredArgs = 2
	if len(cmd.Args) < requiredArgs {
		log.Fatal("Too little arguments provided for given command.")
	} else if len(cmd.Args) > requiredArgs {
		log.Fatal("Too much arguments provided for given command.")
	}

	providedName := cmd.Args[0]
	providedURL := cmd.Args[1]

	user, err := s.Db.GetUser(context.Background(), s.Confg.CurrentUserName)
	if err != nil {
		return err
	}

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
	fmt.Printf("%+v\n", feed)
	return nil
}

// Name: Feeds (no args)