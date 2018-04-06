package file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/jeyem/mogo"
)

var (
	owner            = bson.ObjectIdHex("58237f6bd9d7db2827b9b392")
	uploadableImages = []string{
		"test-image/test.jpg",
		"test-image/test2.jpg",
		"test-image/test3.jpg",
		"test-image/test4.jpg",
		"test-image/test.png",
	}
	shouldServe = []string{
		"dae49be638e75b306cd116bfdb0632d5.jpg",
		"f00a78e13abe880ecdaedbbf5499cb77.png",
		"dae49be638e75b306cd116bfdb0632d5*w250h250.jpg",
		"dae49be638e75b306cd116bfdb0632d5*w150h250.jpg",
		"f00a78e13abe880ecdaedbbf5499cb77*w250h250.jpg",
		"f00a78e13abe880ecdaedbbf5499cb77*w250h150.jpg",
		"dae49be638e75b306cd116bfdb0632d5*w500.jpg",
		"f00a78e13abe880ecdaedbbf5499cb77*h700.jpg",
	}
)

func testPackageinit() {
	cfg := new(Config)
	cfg.ImagePath = "/tmp/test_images"
	config = cfg
	database, err := mogo.Conn("127.0.0.1:27017/test_images")
	if err != nil {
		log.Fatal(err)
	}
	db = database
}

func TestUpload(t *testing.T) {
	testPackageinit()
	ts := httptest.NewServer(http.HandlerFunc(uploadImages))
	defer ts.Close()
	for _, image := range uploadableImages {
		if err := upload(ts.URL, image); err != nil {
			t.Error(err, "on jpeg")
		}
	}
}

func TestServe(t *testing.T) {
	testPackageinit()
	ts := httptest.NewServer(http.HandlerFunc(serveImages))
	defer ts.Close()
	for _, img := range shouldServe {
		url := ts.URL + "/" + img
		if err := getTest(url); err != nil {
			t.Error(err, " ", img)
		}
	}

}

func serveImages(w http.ResponseWriter, r *http.Request) {
	path := GetFile(owner, r.URL.String())
	if _, err := os.Stat(path); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func uploadImages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	upload := New(r, "file", owner)
	_, err := upload.Upload()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"msg": "ok",
	})
	return
}

func upload(url, file string) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add your image file
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	fw, err := w.CreateFormFile("file", file)
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, f); err != nil {
		return
	}
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}

func getTest(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}
	return nil
}
