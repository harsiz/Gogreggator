package main

import (
	"log"                               
	"os"
	"github.com/harsiz/goggregator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	var State = config.State{
		Confg: &conf,
	}
	var Commands = config.Commands{
		CommandMap: make(map[string]func(*config.State, config.Command) error),
	}
	Commands.Register("login", config.HandlerLogin)
	
	// THE VARIABLE NAMES UNDERNEATH *ARE* CONFUSING
	// BARE WITH ME PLEASE


	CommandArguments := os.Args

	if len(CommandArguments) < 2 {
		log.Fatal("Error: Too little arguments.")
	}
	Comm := config.Command{
		Name: CommandArguments[1],
		Args: CommandArguments[2:],
	}

	err = Commands.Run(&State, Comm)
}