package main

import (
	"database/sql"
	"github.com/harsiz/goggregator/internal/config"
	"github.com/harsiz/goggregator/internal/database"
	_ "github.com/lib/pq"
	"log"
	"os"
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

	// REGISTER COMMANDS

	Commands.Register("login", config.HandlerLogin)
	Commands.Register("register", config.HandlerRegister)
	Commands.Register("reset", config.HandlerReset)
	Commands.Register("users", config.HandlerUsers)

	// DB HANDLER

	db, err := sql.Open("postgres", State.Confg.DbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dbQueries := database.New(db)
	State.Db = dbQueries

	// OTHER ETC
	CommandArguments := os.Args

	if len(CommandArguments) < 2 {
		log.Fatal("Error: Too little arguments.")
	}
	Comm := config.Command{
		Name: CommandArguments[1],
		Args: CommandArguments[2:],
	}

	err = Commands.Run(&State, Comm)
	if err != nil {
		log.Fatal(err)
	}
}
