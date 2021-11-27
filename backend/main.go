package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/rs/cors"

	"github.com/gorilla/mux"
)

func handleUploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("+++ handleUploadFile +++")

	vars := mux.Vars(r)
	fmt.Println("name =", vars["name"])

	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("file2upload")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create file
	filePath := path.Join(FOLDER_FOR_UPLOADING, vars["name"])
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Successfully Uploaded File")
	w.WriteHeader(http.StatusOK)
}

var FOLDER_FOR_UPLOADING = "uploadedfiles"

type FileList struct {
	Files []string `json:"files"`
}

func handleGetFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("+++ handleGetFiles +++")

	files, err := ioutil.ReadDir(FOLDER_FOR_UPLOADING)
	if err != nil {
		log.Fatal(err)
	}

	fileArray := []string{}

	for _, f := range files {
		fmt.Println()
		if f.Name() == ".DS_Store" {
			continue
		}
		fileArray = append(fileArray, f.Name())
	}

	w.Header().Set("Content-Type", "application/json")
	fileList := &FileList{Files: fileArray}
	json.NewEncoder(w).Encode(fileList)
}

func handleDeleteFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("+++ handleDeleteFile +++")
	vars := mux.Vars(r)
	fmt.Println("name =", vars["name"])
	filePath := path.Join(FOLDER_FOR_UPLOADING, vars["name"])
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		w.WriteHeader(http.StatusOK)
		return
	}

	e := os.Remove(filePath)

	if e != nil {
		log.Fatal(e)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func setupRoutes() {

	router := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodDelete,
			http.MethodPut,
			http.MethodOptions},
		AllowCredentials: true,
	})

	// List of APIs
	router.HandleFunc("/files/{name}", handleUploadFile).Methods(http.MethodPost)

	router.HandleFunc("/files", handleGetFiles).Methods(http.MethodGet)
	router.HandleFunc("/files/{name}", handleDeleteFile).Methods(http.MethodDelete)

	http.ListenAndServe(":8080", c.Handler(router))
}

func main() {
	fmt.Println("Start Storage Server ...")
	setupRoutes()
}
