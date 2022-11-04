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

## Endpoints

`GET /v1/images`

List metadata for stored images.

`GET /v1/images/<id>`

Get metadata for image with id `<id>`.

`GET /v1/images/<id>/data`

Get image data for image with id `<id>`.

`POST /v1/images`

Upload new image. Request body should be image data.

`PUT /v1/images/<id>`

Update image. Request body should be image data.