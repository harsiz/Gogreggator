package config

import (
	"fmt"
)

func HandlerLogin(s *State, cmd Command) error {
	var requiredArgs = 1
	if len(cmd.Args) < requiredArgs {
		err := fmt.Errorf("%s requires atleast (%d) arguments passed.", cmd.Name, requiredArgs)
		return err
	} else if len(cmd.Args) > requiredArgs {
		err := fmt.Errorf("%s requires no more than (%d) arguments passed.", cmd.Name, requiredArgs)
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