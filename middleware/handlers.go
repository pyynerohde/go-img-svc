package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go-img-svc/models"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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

func GetMetadata(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get imageid from request params and convert type to int
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert type to int. %v", err)
	}

	// Get metadata by calling secondary function getMetadata() with image id
	metadata, err := getMetadata(int64(id))
	if err != nil {
		log.Fatalf("Unable to get metadata by ID. Error: %v", err)
	}
	// Return error 404 if metadata is empty
	if metadata.ImageID == 0 {
		w.WriteHeader(404)
		return
	}

	// Send the response
	err = json.NewEncoder(w).Encode(metadata)
	if err != nil {
		w.WriteHeader(400)
		return
	}
}

func AddImage(w http.ResponseWriter, r *http.Request) {

	// convert request body to bytes
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Save image from base64
	// Extract metadata
	// Add to db

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

func getMetadata(id int64) (models.Image, error) {
	db := createConnection()
	defer db.Close()

	var image models.Image
	sqlStatement := `SELECT * FROM images WHERE imageid=$1`
	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&image.ImageID, &image.Filepath, &image.Filesize, &image.Width, &image.Height, &image.Type, &image.Date)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("Query returned zero rows.")
		return image, nil
	case nil:
		return image, nil
	default:
		log.Fatalf("Unable to scan the row. Error: %v", err)
	}

	return image, err
}
