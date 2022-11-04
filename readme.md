# Image service
A service for storing and retrieving image data, using Go (Golang) and PostgreSQL.

## Setup 
1) Create a db called imagesdb `CREATE DATABASE imagesdb;` and `\c imagesdb` into it. Then add the following table `images`:
```
CREATE TABLE images (
    imageid SERIAL PRIMARY KEY,
    filepath TEXT,
    filesize INT,
    width INT,
    height INT,
    type TEXT,
    date TEXT
    );
```

2) Create a `.env` file containing `POSTGRES_URL="postgres://<user>@localhost:5432/imagesdb?sslmode=disable"` 

3) Change the filepath in const `dir_to_scan` located in `upload.go` to correspond to your machine.


## Endpoints
APIs are found at `http://localhost:8080/api`. 

`GET /v1/images`

List metadata for stored images.

`GET /v1/images/<id>`

Get metadata for image with id `<id>`.

`GET /v1/images/<id>/data`

Get image data for image with id `<id>`.

`POST /v1/images`

Upload new image. Request body should be image data.
* Working with formats PNG, JPEG, GIF.
* Send base64 encoded image in request body.  

`PUT /v1/images/<id>`

Update image. Request body should be image data.

## Shortcuts
* Hardcoded filepaths. 
* Saving images should not all be put in the same folder /img/saved, but instead be saved dynamically.
Eg dynamically save in folders `//ROOT/img/YYYY/MM/DD/HH/MM` depending on the time. Could be more specific with `/SS` etc if needed.