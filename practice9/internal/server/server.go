package server

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"practice9/internal/handlers"
	"practice9/internal/middleware/logging"
)

func Run() {
	r := mux.NewRouter()
	r.HandleFunc("/files", handlers.FileListHandler).Methods("GET")
	r.HandleFunc("/files/{id}", handlers.FileTextHandler).Methods("GET")
	r.HandleFunc("/files/{id}/info", handlers.FileInfoHandler).Methods("GET")
	r.HandleFunc("/files", handlers.UploadFileHandler).Methods("POST")
	r.HandleFunc("/files/{id}", handlers.DeleteFileHandler).Methods("DELETE")
	r.HandleFunc("/files/{id}", handlers.UpdateFileHandler).Methods("PUT")

	http.Handle("/", r)

	logging.Load()

	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Не удалось загрузить файл .env")
	}
	// Чтение переменных окружения
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
