package models

// Schema of the image table
type Image struct {
	ImageID  int64  `json:"imageid"`
	Filepath string `json:"filepath"`
	Filesize int64  `json:"filesize"`
	Width    int64  `json:"width"`
	Height   int64  `json:"height"`
	Type     string `json:"type"`
	Date     string `json:"date"`
}
