package handlers

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"chores-app/internal/db"
	"chores-app/internal/models"
)

type Pages struct {
	DB   *sql.DB
	Tmpl TemplateRenderer
}

type TemplateRenderer interface {
	ExecuteTemplate(w io.Writer, name string, data any) error
}

func (p Pages) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	today := time.Now()
	tasks, err := db.FetchTasks(p.DB, "", today)
	if err != nil {
		http.Error(w, "failed to load tasks", http.StatusInternalServerError)
		return
	}

	data := models.PageData{
		Tasks:   tasks,
		Today:   today.Format("2006-01-02"),
		KidName: "",
		KidList: collectKids(tasks),
	}

	if err := p.Tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("template render: %v", err)
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (p Pages) Kid(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/kid/") {
		http.NotFound(w, r)
		return
	}

	kidName, err := url.PathUnescape(strings.TrimPrefix(r.URL.Path, "/kid/"))
	if err != nil || kidName == "" {
		http.NotFound(w, r)
		return
	}

	today := time.Now()
	tasks, err := db.FetchTasks(p.DB, kidName, today)
	if err != nil {
		http.Error(w, "failed to load tasks", http.StatusInternalServerError)
		return
	}

	data := models.PageData{
		Tasks:   tasks,
		Today:   today.Format("2006-01-02"),
		KidName: kidName,
		KidList: nil,
	}

	if err := p.Tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("template render: %v", err)
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func collectKids(tasks []models.Task) []string {
	seen := make(map[string]struct{})
	for _, t := range tasks {
		seen[t.Kid] = struct{}{}
	}

	kids := make([]string, 0, len(seen))
	for k := range seen {
		kids = append(kids, k)
	}
	sort.Strings(kids)
	return kids
}
