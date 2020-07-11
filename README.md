# FileX

File management for Golang

[![Build Status](http://img.shields.io/travis/prongbang/filex.svg)](https://travis-ci.org/prongbang/filex)
[![Go Report Card](https://goreportcard.com/badge/github.com/prongbang/filex)](https://goreportcard.com/report/github.com/prongbang/filex)


```
go get github.com/prongbang/filex
```

### How to use

- New FileX

```go
fileX := filex.New()
```

- Make Directory

```go
isSuccess := fileX.Mkdir("public/thumbnail")
```

- Create Image from Bytes Array

```go
imgPath := "public/thumbnail/image.jpeg"

imgByte := []byte("mock image byte array")
fileX.CreateImage(imgByte, imgPath)
```

- Create File from Multipart

```go
path := "public/thumbnail"
filename := "image.jpeg"
var file *multipart.FileHeader = mockFile
pathFile, err := fileX.CreateMultipart(path, filename, file)
```

- Delete

```go
path, err := fileX.Delete(imgPath)
```