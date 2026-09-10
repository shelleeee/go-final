package api

import (
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор задачи"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		return
	}

	now := time.Now().UTC()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
			return
		}
	} else {
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}
		if nextDate == "" {
			writeJSON(w, map[string]string{"error": "Не удалось вычислить следующую дату"}, http.StatusBadRequest)
			return
		}

		err = db.UpdateTaskDate(id, nextDate)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
			return
		}
	}

	writeJSON(w, map[string]string{}, http.StatusOK)
}
