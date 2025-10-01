package get_tasks

import "github.com/4udiwe/download-task-service/internal/entity"

type TaskService interface {
	ListTasks() ([]*entity.Task, error)
}
