package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Kredens struct {
	Key   string
	Value string
}

var db *sql.DB

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dbPath := filepath.Join(homeDir, ".kredens.db")
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := initDB(); err != nil {
		log.Fatal(err)
	}

	args := os.Args[1:]
	if len(args) == 0 {
		printHelp()
		return
	}

	switch args[0] {
	case "help":
		printHelp()
	case "list":
		listKredens()
	case "keys":
		listKeys()
	case "vals":
		listVals()
	case "source":
		sourceKredens()
	case "get":
		if len(args) != 2 {
			log.Fatal("Usage: kredens get KEY")
		}
		getKredens(args[1])
	case "set":
		if len(args) != 3 {
			log.Fatal("Usage: kredens set KEY VALUE")
		}
		setKredens(args[1], args[2])
	case "del":
		if len(args) != 2 {
			log.Fatal("Usage: kredens del KEY")
		}
		deleteKredens(args[1])
	default:
		// If only one argument provided, treat it as 'get'
		if len(args) == 1 {
			getKredens(args[0])
			return
		}
		log.Fatalf("Unknown command: %s", args[0])
	}
}

func printHelp() {
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
}

func initDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS kredens (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);`

	_, err := db.Exec(query)
	return err
}

func listKredens() {
	rows, err := db.Query("SELECT key, value FROM kredens")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var k Kredens
		if err := rows.Scan(&k.Key, &k.Value); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s=%s\n", k.Key, k.Value)
	}
	
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func listKeys() {
	rows, err := db.Query("SELECT key FROM kredens")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			log.Fatal(err)
		}
		fmt.Println(key)
	}
	
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func listVals() {
	rows, err := db.Query("SELECT value FROM kredens")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			log.Fatal(err)
		}
		fmt.Println(value)
	}
	
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func sourceKredens() {
	rows, err := db.Query("SELECT key, value FROM kredens")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var k Kredens
		if err := rows.Scan(&k.Key, &k.Value); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("export %s=%s\n", k.Key, k.Value)
	}
	
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func getKredens(key string) {
	var value string
	err := db.QueryRow("SELECT value FROM kredens WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		log.Fatalf("Key not found: %s", key)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(value)
}

func setKredens(key, value string) {
	query := `INSERT OR REPLACE INTO kredens (key, value) VALUES (?, ?)`
	_, err := db.Exec(query, key, value)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteKredens(key string) {
	result, err := db.Exec("DELETE FROM kredens WHERE key = ?", key)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rows == 0 {
		log.Fatalf("Key not found: %s", key)
	}
}
