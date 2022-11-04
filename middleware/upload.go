package middleware

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"go-img-svc/models"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

// Shortcut to have this directory hardcoded. Change this to your own path.
const dir_to_scan string = "/Users/oscarrohde/GolandProjects/go-img-svc/img"

func extractImgInfo() {
	// Search for new files in /img and extract its metadata

	files, _ := ioutil.ReadDir(dir_to_scan)
	for _, imgFile := range files {

		if reader, err := os.Open(filepath.Join(dir_to_scan, imgFile.Name())); err == nil {
			defer reader.Close()
			im, _, err := image.DecodeConfig(reader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", imgFile.Name(), err)
				continue
			}
			fmt.Printf("%s %d %d %d %d\n", imgFile.Name(), im.Width, im.Height, imgFile.Size(), imgFile.ModTime().YearDay())
			fileExtention := strings.SplitAfter(imgFile.Name(), ".")

			// Rename and move image to /img/saved
			oldFilename := dir_to_scan + "/" + imgFile.Name()
			newFilename := dir_to_scan + "/saved/" + strconv.Itoa(imgFile.ModTime().Nanosecond()) + "." + fileExtention[1]
			err = os.Rename(oldFilename, newFilename)
			if err != nil {
				log.Fatalf("Error while renaming image. %v", err)
			}

			// Add above information to Struct + db
			addStructToDb(newFilename, im, imgFile)

		} else {
			fmt.Println("Impossible to open the file:", err)
		}
	}
}

func addStructToDb(filepath string, im image.Config, imgFile fs.FileInfo) *models.Image {
	var img models.Image
	img.Filepath = filepath
	img.Filesize = imgFile.Size()
	img.Width = int64(im.Width)
	img.Height = int64(im.Height)
	img.Type = strings.SplitAfter(imgFile.Name(), ".")[1]
	img.Date = imgFile.ModTime().String()

	// Add image to database with addImage()
	addImage(img)
	return &img
}