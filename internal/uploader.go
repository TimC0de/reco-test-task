package internal

import (
	"os"
	"reco-test-task/internal/common"
)

type Uploader interface {
	UploadNewFile(filename string, data []byte) error
}

type _uploader struct{}

func (u *_uploader) UploadNewFile(filename string, data []byte) error {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}
	fullPath := workingDirectory + common.UPLOAD_PATH + filename + ".json"

	// Return if file already exists
	if _, err := os.Stat(fullPath); os.IsExist(err) {
		return err
	}

	// Insert data into file
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}

func NewUploader() Uploader {
	return &_uploader{}
}
