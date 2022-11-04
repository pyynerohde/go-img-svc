package middleware

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func createConnection() *sql.DB {
	// Load .env file containing the Postgres URL
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Could not load .env file.")
	}

	// Open db connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	// Verify the db connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connection successful!")

	return db
}
