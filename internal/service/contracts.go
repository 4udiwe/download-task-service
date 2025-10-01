package service

import "github.com/4udiwe/download-task-service/internal/entity"

type TaskStorage interface {
	Save(task *entity.Task) error
	Get(id string) (*entity.Task, error)
	List() ([]*entity.Task, error)
	Update(task *entity.Task) error
}
