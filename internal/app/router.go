package app

import (
	"fmt"
	"net/http"

	"github.com/4udiwe/download-task-service/pkg/validator"
	"github.com/labstack/echo/v4"
)

func (app *App) EchoHandler() *echo.Echo {
	if app.echoHandler != nil {
		return app.echoHandler
	}

	handler := echo.New()
	handler.Validator = validator.NewCustomValidator()

	app.configureRouter(handler)

	for _, r := range handler.Routes() {
		fmt.Printf("%s %s\n", r.Method, r.Path)
	}

	app.echoHandler = handler
	return app.echoHandler
}

func (app *App) configureRouter(handler *echo.Echo) {

	tasksGroup := handler.Group("tasks")
	{
		tasksGroup.GET("/:taskID", app.GetTaskByIDHandler().Handle)
		tasksGroup.GET("", app.GetTasksHandler().Handle)
		tasksGroup.POST("", app.PostTaskHandler().Handle)
	}

	handler.GET("/health", func(c echo.Context) error { return c.NoContent(http.StatusOK) })
}
