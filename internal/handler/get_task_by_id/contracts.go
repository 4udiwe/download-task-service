package get_task_by_id

import "github.com/4udiwe/download-task-service/internal/entity"

type TaskService interface {
	GetTask(id string) (*entity.Task, error)
}
