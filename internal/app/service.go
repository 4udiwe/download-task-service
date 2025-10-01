package app

import "github.com/4udiwe/download-task-service/internal/service"

func (app *App) TaskService() *service.TaskService {
	if app.taskService != nil {
		return app.taskService
	}
	app.taskService = service.New(app.TaskStorage(), app.cfg.Download.WorkersCount, nil)
	return app.taskService
}
