package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if repeat == "" {
		return "", nil
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %s", dstart)
	}
	date = date.UTC()
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("пустое правило повторения")
	}

	switch parts[0] {
	case "d":
		if len(parts) < 2 {
			return "", fmt.Errorf("не указано количество дней для правила d")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("неверное значение дня: %s", parts[1])
		}
		if interval > 400 {
			return "", fmt.Errorf("интервал не может превышать 400 дней")
		}
		next := date
		for {
			next = next.AddDate(0, 0, interval)
			if next.After(now) {
				break
			}
		}
		return next.Format(DateFormat), nil

	case "y":
		next := date
		for {
			next = next.AddDate(1, 0, 0)
			if next.After(now) {
				break
			}
		}
		return next.Format(DateFormat), nil

	case "w":
		if len(parts) < 2 {
			return "", fmt.Errorf("не указаны дни недели для правила w")
		}

		weekDays := make(map[time.Weekday]bool)
		dayStrs := strings.Split(parts[1], ",")
		for _, ds := range dayStrs {
			d, err := strconv.Atoi(ds)
			if err != nil {
				return "", fmt.Errorf("неверное значение дня недели: %s", ds)
			}
			if d < 1 || d > 7 {
				return "", fmt.Errorf("день недели должен быть от 1 до 7: %d", d)
			}
			var wd time.Weekday
			if d == 7 {
				wd = time.Sunday
			} else {
				wd = time.Weekday(d)
			}
			weekDays[wd] = true
		}

		next := date.AddDate(0, 0, 1)
		for {
			if weekDays[next.Weekday()] && next.After(now) {
				break
			}
			next = next.AddDate(0, 0, 1)
		}
		return next.Format(DateFormat), nil

	case "m":
		if len(parts) < 2 {
			return "", fmt.Errorf("не указаны дни месяца для правила m")
		}

		dayStrs := strings.Split(parts[1], ",")
		days := make(map[int]bool)
		for _, ds := range dayStrs {
			d, err := strconv.Atoi(ds)
			if err != nil {
				return "", fmt.Errorf("неверное значение дня месяца: %s", ds)
			}
			if d < -2 || d == 0 || d > 31 {
				return "", fmt.Errorf("день месяца должен быть от 1 до 31, -1 или -2: %d", d)
			}
			days[d] = true
		}

		months := make(map[int]bool)
		if len(parts) >= 3 {
			monthStrs := strings.Split(parts[2], ",")
			for _, ms := range monthStrs {
				m, err := strconv.Atoi(ms)
				if err != nil {
					return "", fmt.Errorf("неверное значение месяца: %s", ms)
				}
				if m < 1 || m > 12 {
					return "", fmt.Errorf("месяц должен быть от 1 до 12: %d", m)
				}
				months[m] = true
			}
		} else {
			for m := 1; m <= 12; m++ {
				months[m] = true
			}
		}

		next := date
		for {
			next = next.AddDate(0, 0, 1)

			if next.Sub(now) > 400*24*time.Hour {
				return "", fmt.Errorf("не найдена подходящая дата в пределах 400 дней")
			}

			if !months[int(next.Month())] {
				continue
			}

			year, month, day := next.Date()
			lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

			matched := false
			for d := range days {
				if d >= 1 && d <= 31 && day == d {
					matched = true
					break
				}
				if d == -1 && day == lastDay {
					matched = true
					break
				}
				if d == -2 && day == lastDay-1 {
					matched = true
					break
				}
			}

			if matched && next.After(now) {
				return next.Format(DateFormat), nil
			}
		}

	default:
		return "", fmt.Errorf("неподдерживаемый формат повторения: %s", parts[0])
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "Неверный формат параметра now", http.StatusBadRequest)
			return
		}
	}

	if dateStr == "" {
		http.Error(w, "Не указан параметр date", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}
