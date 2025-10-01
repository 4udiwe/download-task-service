package app

import "github.com/4udiwe/download-task-service/internal/storage"

func (app *App) TaskStorage() *storage.Storage {
	if app.taskStorage != nil {
		return app.taskStorage
	}
	app.taskStorage = storage.New(app.cfg.Download.Path)
	return app.taskStorage
}
