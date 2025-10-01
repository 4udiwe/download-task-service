package app

import (
	"github.com/4udiwe/download-task-service/internal/handler"
	"github.com/4udiwe/download-task-service/internal/handler/get_task_by_id"
	"github.com/4udiwe/download-task-service/internal/handler/get_tasks"
	post_task "github.com/4udiwe/download-task-service/internal/handler/post_tast"
)

func (app *App) GetTaskByIDHandler() handler.Handler {
	if app.getTaskByIDHandler != nil {
		return app.getTaskByIDHandler
	}
	app.getTaskByIDHandler = get_task_by_id.New(app.TaskService())
	return app.getTaskByIDHandler
}

func (app *App) GetTasksHandler() handler.Handler {
	if app.getTasksHandler != nil {
		return app.getTasksHandler
	}
	app.getTasksHandler = get_tasks.New(app.TaskService())
	return app.getTasksHandler
}

func (app *App) PostTaskHandler() handler.Handler {
	if app.postTaskHandler != nil {
		return app.postTaskHandler
	}
	app.postTaskHandler = post_task.New(app.TaskService())
	return app.postTaskHandler
}
