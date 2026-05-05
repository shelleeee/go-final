package db

import (
	"database/sql"
	"fmt"
	"strings"
)

func AddTask(task *Task) (int64, error) {
	query := `
        INSERT INTO scheduler (date, title, comment, repeat)
        VALUES (:date, :title, :comment, :repeat)
    `

	result, err := DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func Tasks(limit int, search string) ([]*Task, error) {
	var query string
	var args []interface{}

	isDate := false
	var searchDate string
	if search != "" {
		parts := strings.Split(search, ".")
		if len(parts) == 3 {
			searchDate = parts[2] + parts[1] + parts[0]
			isDate = true
		}
	}

	if isDate {
		query = `
            SELECT id, date, title, comment, repeat
            FROM scheduler
            WHERE date = :date
            ORDER BY date ASC
            LIMIT :limit
        `
		args = []interface{}{
			sql.Named("date", searchDate),
			sql.Named("limit", limit),
		}
	} else if search != "" {
		query = `
            SELECT id, date, title, comment, repeat
            FROM scheduler
            WHERE title LIKE :search OR comment LIKE :search
            ORDER BY date ASC
            LIMIT :limit
        `
		searchPattern := "%" + search + "%"
		args = []interface{}{
			sql.Named("search", searchPattern),
			sql.Named("limit", limit),
		}
	} else {
		query = `
            SELECT id, date, title, comment, repeat
            FROM scheduler
            ORDER BY date ASC
            LIMIT :limit
        `
		args = []interface{}{
			sql.Named("limit", limit),
		}
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var task Task

	query := `
        SELECT id, date, title, comment, repeat
        FROM scheduler
        WHERE id = :id
    `

	row := DB.QueryRow(query, sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `
        UPDATE scheduler
        SET date = :date, title = :title, comment = :comment, repeat = :repeat
        WHERE id = :id
    `

	result, err := DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача с ID %s не найдена", task.ID)
	}

	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = :id`

	result, err := DB.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача с ID %s не найдена", id)
	}

	return nil
}

func UpdateTaskDate(id string, date string) error {
	query := `UPDATE scheduler SET date = :date WHERE id = :id`

	result, err := DB.Exec(query,
		sql.Named("date", date),
		sql.Named("id", id))
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача с ID %s не найдена", id)
	}

	return nil
}
