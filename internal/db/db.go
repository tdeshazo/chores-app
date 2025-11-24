package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// Open returns a SQLite connection configured for this app and ensures schema + seed data exist.
func Open(dbPath string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(1)

	if err := initSchema(conn); err != nil {
		conn.Close()
		return nil, err
	}
	if err := seedTasks(conn); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func initSchema(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			kid TEXT NOT NULL,
			title TEXT NOT NULL,
			sort_order INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS task_log (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			task_id INTEGER NOT NULL,
			date TEXT NOT NULL,
			status TEXT NOT NULL,
			UNIQUE(task_id, date),
			FOREIGN KEY(task_id) REFERENCES tasks(id)
		)`,
	}

	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}

func seedTasks(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count); err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	type seedTask struct {
		kid       string
		title     string
		sortOrder int
	}

	sample := []seedTask{
		{kid: "Griffin", title: "Make bed", sortOrder: 1},
		{kid: "Griffin", title: "Brush teeth", sortOrder: 2},
		{kid: "Griffin", title: "Feed the cat", sortOrder: 3},
		{kid: "Garreth", title: "Put toys away", sortOrder: 1},
		{kid: "Garreth", title: "Set the table", sortOrder: 2},
	}

	for _, t := range sample {
		if _, err := db.Exec(
			"INSERT INTO tasks (kid, title, sort_order) VALUES (?, ?, ?)",
			t.kid, t.title, t.sortOrder,
		); err != nil {
			return err
		}
	}

	return nil
}
