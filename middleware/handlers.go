package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"go-img-svc/models"
	"log"
	"net/http"
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

func GetAllMetadata(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get all the metadata in the database using getAllMetadata()
	metadata, err := getAllMetadata()
	if err != nil {
		log.Fatalf("Unable to get all metadata. Error: %v", err)
	}

	// Send all the metadata as a response
	err = json.NewEncoder(w).Encode(metadata)
	if err != nil {
		w.WriteHeader(400)
		return
	}
}

/* ####################################################################################### */
/* ---------------------------------- Handler functions ---------------------------------- */
/* ####################################################################################### */

func getAllMetadata() ([]models.Image, error) {
	db := createConnection()
	defer db.Close()

	var images []models.Image
	sqlStatement := `SELECT * FROM images`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. Error: %v", err)
	}
	defer rows.Close()

	// iterate the rows
	for rows.Next() {
		var image models.Image

		err = rows.Scan(&image.ImageID, &image.Filepath, &image.Filesize, &image.Width, &image.Height, &image.Type, &image.Date)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		images = append(images, image)
	}

	return images, err
}
