package main

import (
	"excelize/api"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	// middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Send a POST request to /excel with a JSON file to convert it to Excel. (multipart/form-data with key 'file')"))
	})

	r.Get("/dummy", api.DummyHandler)
	r.Post("/excel", api.ExcelHandler)

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
