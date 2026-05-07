package main

import (
	"log"
	"os"

	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	if err := db.Init(dbFile); err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	defer db.Close()

	if err := server.Run(); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}
