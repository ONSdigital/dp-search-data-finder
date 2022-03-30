package models

import (
	"time"
)

// Task represents a metadata model and json representation for API
type Task struct {
	JobID             string     `json:"job_id"`
	LastUpdated       time.Time  `json:"last_updated"`
	Links             *TaskLinks `json:"links"`
	NumberOfDocuments int        `json:"number_of_documents"`
	TaskName          string     `json:"task_name"`
}

// TaskLinks is a type that contains links to the endpoints for returning a specific task (self), and the job that it is part of (job), respectively.
type TaskLinks struct {
	Self string `json:"self"`
	Job  string `json:"job"`
}
