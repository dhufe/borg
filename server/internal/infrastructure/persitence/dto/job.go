package dto

import "time"

type (
	JobStatus string
	JobType   string
)

const (
	StatusPending   JobStatus = "Pending"
	StatusRunning   JobStatus = "Running"
	StatusCompleted JobStatus = "Completed"
	StatusFailed    JobStatus = "Failed"
)

type Job struct {
	ID          uint       `gorm:"primaryKey"      json:"id"`
	FileName    string     `gorm:"not null"        json:"filename"`
	Status      JobStatus  `gorm:"default:Pending" json:"status"`
	StartedAt   *time.Time `                       json:"started_at"`
	CompletedAt *time.Time `                       json:"completed_at"`
	Output      string     `gorm:"type:text"       json:"output"` // Ausgabe des Programms
	Error       string     `gorm:"type:text"       json:"error"`  // Fehlerausgabe
	CreatedAt   time.Time  `gorm:"autoCreateTime"  json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"  json:"updated_at"`
}

func AllJobStatuses() []JobStatus {
	return []JobStatus{
		StatusPending,
		StatusRunning,
		StatusCompleted,
		StatusFailed,
	}
}
