package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
)

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func Test_handleGetFiles(t *testing.T) {

	request, _ := http.NewRequest("GET", "/files", nil)
	response := httptest.NewRecorder()
	handleGetFiles(response, request)

	if response.Code != 200 {
		t.Fatalf("Non-expected status code %v:\n\tbody: %v", "200", response.Code)
	} else {
		t.Log("Reponse:", response)
	}
}

func Test_handleDeleteFile(t *testing.T) {

	copyFile("mylogo.png", "./uploadedfiles/mylogo.png")

	request, _ := http.NewRequest("GET", "/files/mylogo.png", nil)
	response := httptest.NewRecorder()

	request = mux.SetURLVars(request, map[string]string{
		"name": "mylogo.png",
	})

	handleDeleteFile(response, request)

	if response.Code != 200 {
		t.Fatalf("Non-expected status code %v:\n\tbody: %v", "200", response.Code)
	} else {
		t.Log("Reponse:", response)
	}
}

func Test_handleUploadFile(t *testing.T) {

	var err error

	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	const fileName = "main.go"

	//prepare the reader instances to encode
	values := map[string]io.Reader{
		"file2upload": mustOpen(fileName), // lets assume its this file
	}

	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	t.Log("Ready to upload file:", fileName)
	url := "/files/" + fileName
	t.Log("URL =", url)

	// Now that you have a form, you can submit it to your handler.
	request, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}

	// Don't forget to set the content type, this will contain the boundary.
	request.Header.Set("Content-Type", w.FormDataContentType())

	response := httptest.NewRecorder()
	handleUploadFile(response, request)

	t.Log("Reponse:", response)

	if response.Code != 200 {
		t.Fatalf("Non-expected status code %v:\n\tbody: %v", "200", response.Code)
	}
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
