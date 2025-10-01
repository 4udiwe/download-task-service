package entity

import "time"

type TaskStatus string

const (
	TaskPending TaskStatus = "pending"
	TaskRunning TaskStatus = "in_progress"
	TaskDone    TaskStatus = "done"
	TaskFailed  TaskStatus = "failed"
)

type Task struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	Status    TaskStatus `json:"status"`
	Files     []*File    `json:"files"`
}
