package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go_final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Неверный формат JSON"}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	today := now.Format(DateFormat)

	if task.Date == "" {
		task.Date = today
	}

	dateTime, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Неверный формат даты"}, http.StatusBadRequest)
		return
	}
	dateTime = dateTime.UTC()
	dateTime = time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 0, 0, 0, 0, time.UTC)

	if dateTime.Before(now) {
		if task.Repeat == "" {
			task.Date = today
		} else {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
				return
			}
			if nextDate != "" {
				task.Date = nextDate
			}
		}
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка добавления задачи"}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)}, http.StatusCreated)
}
