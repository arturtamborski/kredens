package main

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func run() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dbPath := filepath.Join(homeDir, ".kredens.db")
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := initDB(); err != nil {
		return err
	}

	return handleArgs()
}

func handleArgs() error {
	args := os.Args[1:]
	if len(args) == 0 {
		return helpCmd()
	}

	switch args[0] {
	case "help":
		return helpCmd()
	case "list":
		return listCmd(ListKeys | ListValues)
	case "keys":
		return listCmd(ListKeys)
	case "vals":
		return listCmd(ListValues)
	case "values":
		return listCmd(ListValues)
	case "source":
		return sourceCmd()

	case "get":
		if len(args) != 2 {
			return errors.New("usage: kredens get KEY")
		}
		return getCmd(args[1])

	case "set":
		if len(args) != 3 {
			return errors.New("usage: kredens set KEY VALUE")
		}
		return setCmd(args[1], args[2])

	case "del":
		if len(args) != 2 {
			return errors.New("usage: kredens del KEY")
		}
		return delCmd(args[1])

	default:
		return errors.New("unknown command")
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
