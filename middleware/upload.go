package middleware

import (
	"bytes"
	"encoding/base64"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

func saveOnDisk(data string) bool {
	// Configure name of file
	idx := strings.Index(data, ";base64,")
	if idx < 0 {
		panic("InvalidImage")
	}
	ImageType := data[11:idx]
	log.Println(ImageType)

	unbased, err := base64.StdEncoding.DecodeString(data[idx+8:])
	if err != nil {
		log.Println("Cannot decode b64")
		return false
	}
	r := bytes.NewReader(unbased)

	switch ImageType {
	case "png":
		im, err := png.Decode(r)
		if err != nil {
			log.Println("PNG base64 file is broken and can't be read.")
			return false
		}

		f, err := os.OpenFile("img/new.png", os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			log.Println("Cannot open file")
			return false
		}
		png.Encode(f, im)

	case "jpeg":
		im, err := jpeg.Decode(r)
		if err != nil {
			log.Println("JPEG base64 file is broken and can't be read.")
			return false
		}

		f, err := os.OpenFile("img/new.jpeg", os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			log.Println("Cannot open file")
			return false
		}
		jpeg.Encode(f, im, nil)

	case "gif":
		im, err := gif.Decode(r)
		if err != nil {
			log.Println("GIF base64 file is broken and can't be read.")
			return false
		}

		f, err := os.OpenFile("img/new.gif", os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			log.Println("Cannot open file")
			return false
		}
		gif.Encode(f, im, nil)
	}
	return true
}
