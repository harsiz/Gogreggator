package config

import (
	"context"
	"github.com/harsiz/goggregator/internal/database"
)

func MiddlewareLoggedIn(
	handler func(*State, Command, database.User) error,
) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		u := s.Confg.CurrentUserName
		user, err := s.Db.GetUser(context.Background(), u)
		if err != nil {
			return err
		}
		err = handler(s, cmd, user)
		if err != nil {
			return err
		}
		return nil
	}
}
