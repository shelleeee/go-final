package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

type tasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(TasksLimit, search)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка получения задач"}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, tasksResponse{Tasks: tasks}, http.StatusOK)
}
