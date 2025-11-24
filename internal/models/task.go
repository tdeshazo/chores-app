package models

// Task represents a chore and its status for rendering.
type Task struct {
	ID        int
	Kid       string
	Title     string
	SortOrder int
	Status    string
}

// PageData contains the fields needed to render the main template.
type PageData struct {
	Tasks   []Task
	Today   string
	KidName string
	KidList []string
}
