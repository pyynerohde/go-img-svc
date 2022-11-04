# Image service
A service for storing and retrieving image data.

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