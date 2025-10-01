package post_task

import "github.com/4udiwe/download-task-service/internal/entity"

type TaskService interface {
	CreateTask(urls []string) (*entity.Task, error)
}
