package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type API struct {
	DB *sql.DB
}

func (a API) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		TaskID int    `json:"task_id"`
		Status string `json:"status"`
	}

	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	defer r.Body.Close()
	var body payload
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json payload")
		return
	}

	if body.TaskID == 0 {
		writeJSONError(w, http.StatusBadRequest, "missing task_id")
		return
	}

	if !isValidStatus(body.Status) {
		writeJSONError(w, http.StatusBadRequest, "invalid status")
		return
	}

	today := time.Now().Format("2006-01-02")
	_, err := a.DB.Exec(
		`INSERT INTO task_log (task_id, date, status)
			 VALUES (?, ?, ?)
			 ON CONFLICT(task_id, date)
			 DO UPDATE SET status = excluded.status`,
		body.TaskID, today, body.Status,
	)
	if err != nil {
		log.Printf("update status: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "database error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func isValidStatus(status string) bool {
	switch status {
	case "pending", "done", "skipped":
		return true
	default:
		return false
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("write json: %v", err)
	}
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
