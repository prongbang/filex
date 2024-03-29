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

// FileX is a interface
type FileX interface {
	CreateMultipart(path string, filename string, file *multipart.FileHeader) (string, error)
	CreateImage(imgByte []byte, path string) (string, error)
	CreateFile(path string, fileName string, data string) (string, error)
	Delete(path string) (string, error)
	DeleteDir(path string) error
	Mkdir(dirName string) bool
	Read(path string) []byte
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
	defer func(src multipart.File) {
		_ = src.Close()
	}(src)

	ext := filepath.Ext(file.Filename)
	filePath := fmt.Sprintf("%s%s", filename, ext)

	f.Mkdir(path)
	pathFile := fmt.Sprintf("%s/%s", path, filePath)
	dst, err := os.Create(pathFile)
	if err != nil {
		return "", err
	}
	defer func(dst *os.File) {
		_ = dst.Close()
	}(dst)

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

func (f *fileX) DeleteDir(path string) error {
	return os.RemoveAll(path)
}

func (f *fileX) CreateImage(imgByte []byte, path string) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		log.Println(err)
		return "", err
	}

	out, _ := os.Create(path)
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	var opts jpeg.Options
	opts.Quality = 100

	err = jpeg.Encode(out, img, &opts)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return "", nil
}

func (f *fileX) CreateFile(path string, fileName string, data string) (string, error) {
	f.Mkdir(path)
	pathFile := fmt.Sprintf("%s/%s", path, fileName)
	file, err := os.Create(pathFile)
	if err != nil {
		fmt.Println(err)
	}

	_, err = file.WriteString(data)
	if err != nil {
		_ = file.Close()
		return "", err
	}

	err = file.Close()
	if err != nil {
		return "", err
	}

	return pathFile, err
}

func (f *fileX) Read(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		return []byte{}
	}
	return data
}

// New a instance
func New() FileX {
	return &fileX{}
}
