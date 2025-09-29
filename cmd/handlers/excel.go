package handlers

import (
	"excelize/cmd/utils"
	"fmt"
	"net/http"
	"os"
	"time"
)

// ExcelHandler will receive a json file and pass it to the GenerateExcelFromJson to process the Excel generation
func ExcelHandler(w http.ResponseWriter, r *http.Request) {
	filepath, err := utils.UploadFile(w, r)
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	// pass the file path to the cmd package
	excelFile, err := utils.GenerateExcelFromJson(filepath)
	if err != nil {
		http.Error(w, "Failed to generate Excel", http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := os.Remove(excelFile); err != nil {
			fmt.Println("Error removing file:", err)
		}
	}()

	filename := fmt.Sprintf("excel_%d.xlsx", time.Now().Unix())

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	http.ServeFile(w, r, excelFile)
}
