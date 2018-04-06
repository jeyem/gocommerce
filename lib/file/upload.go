package file

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var (
	ErrorNotValidFile = errors.New("file: not valid")
)

type Data struct {
	req   *http.Request
	field string
	owner bson.ObjectId
	files []File
}

func New(r *http.Request, field string, owner bson.ObjectId) *Data {
	return &Data{
		req:   r,
		field: field,
		owner: owner,
	}
}

func (d *Data) Upload() ([]File, error) {
	if err := d.req.ParseMultipartForm(MaxBodySize); err != nil {
		return nil, err
	}
	mpf := d.req.MultipartForm
	fileHeaders, ok := mpf.File[d.field]
	if !ok {
		return d.files, errors.New("upload field not found")
	}
	for _, header := range fileHeaders {
		if err := d.getInput(header); err != nil {
			return d.files, errors.New(header.Filename + " " + err.Error())
		}
	}
	return d.files, nil
}

func (d *Data) getInput(f *multipart.FileHeader) error {
	if f.Size > MaxFileSize {
		return ErrorNotValidFile
	}
	f.Filename = strings.ToLower(f.Filename)
	ext := filepath.Ext(f.Filename)
	ext = strings.Replace(ext, ".", "", -1)
	format, ok := Format[ext]
	if !ok {
		return ErrorNotValidFile
	}
	src, err := f.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	checksum := getMD5Checksum(src)
	subDirectory := filepath.Join(d.owner.Hex(), time.Now().Format("200601"))
	filename := checksum + "." + ext
	path := filepath.Join(subDirectory, filename)
	fullDirectoryPath := filepath.Join(config.ImagePath, subDirectory)
	if _, err := os.Stat(fullDirectoryPath); os.IsNotExist(err) {
		os.MkdirAll(fullDirectoryPath, 0777)
	}
	fullPath := filepath.Join(config.ImagePath, path)
	src.Seek(0, 0)
	dst, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer dst.Close()
	io.Copy(dst, src)
	file := File{
		Owner:    d.owner,
		Name:     filename,
		Path:     path,
		Format:   format,
		CheckSum: checksum,
	}
	if file.Save(); err != nil {
		return err
	}
	d.files = append(d.files, file)
	return nil
}

func getMD5Checksum(f io.Reader) string {
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
