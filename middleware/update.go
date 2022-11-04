package middleware

import (
	"fmt"
	"go-img-svc/models"
	"image"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Shortcut. The code does not remove the saved image on disk.
// It only updates the DB with the corresponding ID, and adds a new image to disk.

func updateMetadata(id int64) {
	// This function is pretty much the same as extractImgInfo(). Shortcut, instead of making it dynamic.
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
			updateDb(id, newFilename, im, imgFile)

		} else {
			fmt.Println("Impossible to open the file:", err)
		}
	}
}

func updateDb(id int64, filepath string, im image.Config, imgFile fs.FileInfo) *models.Image {
	var img models.Image
	img.Filepath = filepath
	img.Filesize = imgFile.Size()
	img.Width = int64(im.Width)
	img.Height = int64(im.Height)
	img.Type = strings.SplitAfter(imgFile.Name(), ".")[1]
	img.Date = imgFile.ModTime().String()

	updateImage(id, img)
	return &img
}
