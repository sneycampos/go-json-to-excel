package internal

import (
	"net/http"
	"os"
)

func UploadFile(w http.ResponseWriter, r *http.Request) (filepath string, err error) {
	defer r.Body.Close()
	r.ParseForm()
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Size > 500*1024*1024 {
		http.Error(w, "File size exceeds 500MB", http.StatusBadRequest)
		return
	}

	workDir, err := os.Getwd()
	if err != nil {
		http.Error(w, "Failed to get working directory", http.StatusInternalServerError)
		return
	}

	// creates the temp dir if not exists
	if _, err := os.Stat(workDir + "/temp"); os.IsNotExist(err) {
		err := os.Mkdir(workDir+"/temp", os.ModePerm)
		if err != nil {
			http.Error(w, "Failed to create temp directory", http.StatusInternalServerError)
		}
	}

	// upload the file to a temp location
	tempFilePath := workDir + "/temp/" + header.Filename
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	_, err = tempFile.ReadFrom(file)
	if err != nil {
		http.Error(w, "Failed to save temp file", http.StatusInternalServerError)
		return
	}

	return tempFilePath, nil
}
