package db

import (
	"database/sql"
	"time"

	"chores-app/internal/models"
)

// FetchTasks returns tasks with today's status, optionally filtered by kid.
func FetchTasks(db *sql.DB, kid string, today time.Time) ([]models.Task, error) {
	var args []any
	todayStr := today.Format("2006-01-02")

	query := `
		SELECT
			t.id,
			t.kid,
			t.title,
			t.sort_order,
			COALESCE(l.status, 'pending') AS status
		FROM tasks t
		LEFT JOIN task_log l
			ON l.task_id = t.id AND l.date = ?
	`
	args = append(args, todayStr)

	if kid != "" {
		query += " WHERE t.kid = ?"
		args = append(args, kid)
	}

	query += " ORDER BY t.kid, t.sort_order"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Kid, &t.Title, &t.SortOrder, &t.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
