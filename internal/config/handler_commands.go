package config

import (
	"fmt"
	"log"
)

func HandlerLogin(s *State, cmd Command) error {
	var requiredArgs = 1
	if len(cmd.Args) < requiredArgs {
		log.Fatal("Too little arguments provided for given command.")
	} else if len(cmd.Args) > requiredArgs {
		log.Fatal("Too much arguments provided for given command.")
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
