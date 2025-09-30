package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"
)

const (
	ListKeys   = 1
	ListValues = 2
)

type Entry struct {
	Key   string
	Value string
}

var db *sql.DB

func initDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS kredens (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);`

	_, err := db.Exec(query)
	return err
}

func listCmd(mode int) error {
	rows, err := db.Query("SELECT key, value FROM kredens ORDER BY key")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.Key, &e.Value); err != nil {
			return err
		}

		if mode == ListKeys {
			fmt.Println(e.Key)
		} else if mode == ListValues {
			fmt.Println(e.Value)
		} else if mode == ListKeys|ListValues {
			fmt.Printf("%s=%s\n", e.Key, e.Value)
		} else {
			return errors.New("unknown mode")
		}
	}

	return rows.Err()
}

func getCmd(key string) error {
	var value string
	err := db.QueryRow("SELECT value FROM kredens WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return fmt.Errorf("key not found: %s", key)
	}

	if err != nil {
		return err
	}

	fmt.Println(value)

	return nil
}

func setCmd(key, value string) error {
	query := `INSERT OR REPLACE INTO kredens (key, value) VALUES (?, ?)`
	_, err := db.Exec(query, key, value)
	if err != nil {
		return err
	}

	return nil
}

func delCmd(key string) error {
	result, err := db.Exec("DELETE FROM kredens WHERE key = ?", key)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("key not found: %s", key)
	}

	return nil
}

func sourceCmd() error {
	rows, err := db.Query("SELECT key, value FROM kredens ORDER BY key")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.Key, &e.Value); err != nil {
			return err
		}
		fmt.Printf("export %s=%s\n", e.Key, e.Value)
	}

	return rows.Err()
}

func helpCmd() error {
	fmt.Println("Usage: kredens [command] [args...]")
	fmt.Println("\nCommands:")
	fmt.Println("  list          List all credentials")
	fmt.Println("  keys          List all keys")
	fmt.Println("  vals          List all values")
	fmt.Println("  get KEY       Show value for KEY")
	fmt.Println("  set KEY VAL   Store KEY with value VAL")
	fmt.Println("  del KEY       Delete KEY")
	fmt.Println("  source        Output credentials as export statements")
	fmt.Println("  help          Show this help message")

	return nil
}
