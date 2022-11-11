package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
		w.WriteHeader(400)
		return
		//log.Fatalf("Unable to convert type to int. %v", err)
	}

	// Get metadata by calling secondary function getMetadata() with image id
	metadata, err := getMetadata(int64(id))
	if err != nil {
		w.WriteHeader(400)
		return
		//log.Fatalf("Unable to get metadata by ID. Error: %v", err)
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

func GetImage(w http.ResponseWriter, r *http.Request) {
	// Get imageid from request params and convert type to int
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(400)
		return
		//log.Fatalf("Unable to convert type to int. %v", err)
	}

	// Get filepath by calling getImage() with image id
	img, err := getImage(int64(id))
	if err != nil {
		w.WriteHeader(400)
		return
		//log.Fatalf("Unable to get metadata by ID. Error: %v", err)
	}
	if img.ImageID == 0 {
		w.WriteHeader(404)
		return
	}

	// Convert image to base64
	path := img.Filepath
	base64img := imgToBase64(path)

	// Send the response
	err = json.NewEncoder(w).Encode(base64img)
	if err != nil {
		w.WriteHeader(400)
		return
	}
}

func AddImage(w http.ResponseWriter, r *http.Request) {

	// Request body should be the whole base64 encoding of an image.
	// convert request body to bytes
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Save image from base64
	saveSuccess := saveOnDisk(string(bytes))
	if saveSuccess != true {
		w.WriteHeader(400)
		return
	}

	// Extract metadata and add to db
	extractImgInfo()
}

func UpdateImage(w http.ResponseWriter, r *http.Request) {
	// Get imageid from request params and convert type to int
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(400)
		return
		//log.Fatalf("Unable to convert type to int. %v", err)
	}

	// Verify that the ID exists already and can be updated
	var image models.Image
	rowsToUpdate := updateImage(int64(id), image)
	if rowsToUpdate == 0 {
		w.WriteHeader(404)
		return
	}

	// convert request body to bytes, then to string for saveOnDisk() to read
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Save updated image on disk
	saveSuccess := saveOnDisk(string(bytes))
	if saveSuccess != true {
		w.WriteHeader(400)
	}

	// Update metadata in db, using updateImage
	updateMetadata(int64(id))
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

func getImage(id int64) (models.Image, error) {
	db := createConnection()
	defer db.Close()

	// Query to get the filepath of the image, using the ID (taken from GetImage)
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

func addImage(image models.Image) int64 {
	db := createConnection()
	defer db.Close()

	var id int64
	sqlStatement := `INSERT INTO images (filepath, filesize, width, height, type, date) VALUES ($1, $2, $3, $4, $5, $6) RETURNING imageid`
	err := db.QueryRow(sqlStatement, image.Filepath, image.Filesize, image.Width, image.Height, image.Type, image.Date).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. Error: %v", err)
	}
	fmt.Printf("Image added successully: %v", id)
	return id
}

func updateImage(id int64, image models.Image) int64 {
	// Shortcut. Does not delete old entry. Could be done by getting filepath and use os.Remove(filepath).

	db := createConnection()
	defer db.Close()

	sqlStatement := `UPDATE images SET filepath=$2, filesize=$3, width=$4, height=$5, type=$6, date=$7 WHERE imageid=$1`
	res, err := db.Exec(sqlStatement, id, image.Filepath, image.Filesize, image.Width, image.Height, image.Type, image.Date)

	if err != nil {
		log.Fatalf("Unable to execute the query. Error: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while performing RowsAffected(). Error: %v", err)
	}
	if rowsAffected == 0 {
		fmt.Println("Query returned zero rows.")
		return rowsAffected
	}
	fmt.Printf("Image was updated successfully. Images updated: %v", rowsAffected)
	return rowsAffected
}
