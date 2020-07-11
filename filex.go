package filex

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileX interface {
	CreateMultipart(path string, filename string, file *multipart.FileHeader) (string, error)
	CreateImage(imgByte []byte, path string) (string, error)
	Delete(path string) (string, error)
	Mkdir(dirName string) bool
}

type fileX struct {
}

func (f *fileX) Mkdir(dirName string) bool {
	src, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirName, 0755)
		if errDir != nil {
			log.Println(errDir)
			return false
		}
		return true
	}
	if src.Mode().IsRegular() {
		log.Println(dirName, "Already exist.")
		return false
	}
	return false
}

func (f *fileX) CreateMultipart(path string, filename string, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	filePath := fmt.Sprintf("%s%s", filename, ext)

	f.Mkdir(path)
	pathFile := fmt.Sprintf("%s/%s", path, filePath)
	dst, err := os.Create(pathFile)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	return pathFile, err
}

func (f *fileX) Delete(path string) (string, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", err
	}
	if err = os.Remove(path); err != nil {
		return "", err
	}
	return path, nil
}

func (f *fileX) CreateImage(imgByte []byte, path string) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		log.Println(err)
		return "", err
	}

	out, _ := os.Create(path)
	defer out.Close()

	var opts jpeg.Options
	opts.Quality = 100

	err = jpeg.Encode(out, img, &opts)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return "", nil
}

func New() FileX {
	return &fileX{}
}
