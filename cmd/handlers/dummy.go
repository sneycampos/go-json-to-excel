package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// DummyHandler generates a large JSON file with 500,000 rows and 10 columns and serves it as a download
func DummyHandler(w http.ResponseWriter, r *http.Request) {
	jsonPath := filepath.Join("temp", "dummy.json")
	file, err := os.Create(jsonPath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	file.Write([]byte("{"))
	for i := 1; i <= 500000; i++ {
		file.Write([]byte(`"` + strconv.Itoa(i) + `": [`))
		file.Write([]byte("{"))
		for j := 0; j < 10; j++ {
			colName := string(rune('A' + j))
			var value string
			if i == 1 {
				value = `"Header ` + colName + `"`
			} else {
				value = `"Data ` + colName + ` Row  ` + strconv.Itoa(i) + `"`
			}
			file.Write([]byte(`"` + colName + `": ` + value))
			if j < 9 {
				file.Write([]byte(",")) // Only add comma if not last column
			}
		}
		file.Write([]byte("}"))
		file.Write([]byte("]"))
		if i < 500000 {
			file.Write([]byte(",")) // Only add comma if not last row
		}
	}
	file.Write([]byte("}"))
	file.Sync()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", jsonPath))
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, jsonPath)
}
