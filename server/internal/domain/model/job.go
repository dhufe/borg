package model

import "time"

type (
	JobStatus string
)

type Job struct {
	ID          uint       `json:"id"`
	FileName    string     `json:"filename"`
	Status      JobStatus  `json:"status"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Output      string     `json:"output"` // Ausgabe des Programms
	Error       string     `json:"error"`  // Fehlerausgabe
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
