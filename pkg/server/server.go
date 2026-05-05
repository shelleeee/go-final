package server

import (
	"log"
	"net/http"
	"os"

	"go_final_project/pkg/api"
)

func Run() error {
	api.Init()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	webDir := "./web"
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	log.Printf("Сервер запущен на порту %s", port)
	return http.ListenAndServe(":"+port, nil)
}
