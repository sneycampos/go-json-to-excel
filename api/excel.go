package api

import (
	"excelize/internal"
	"fmt"
	"net/http"
	"os"
	"time"
)

// ExcelHandler will receive a json file and pass it to the GenerateExcelFromJson to process the Excel generation
func ExcelHandler(w http.ResponseWriter, r *http.Request) {
	filepath, err := internal.UploadFile(w, r)
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	// pass the file path to the cmd package
	excelFile, err := internal.GenerateExcelFromJson(filepath)
	if err != nil {
		http.Error(w, "Failed to generate Excel", http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("excel_%d.xlsx", time.Now().Unix())

	defer os.Remove(excelFile)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	http.ServeFile(w, r, excelFile)
}
