package get_tasks

import (
	"net/http"

	h "github.com/4udiwe/download-task-service/internal/handler"
	"github.com/4udiwe/download-task-service/internal/handler/decorator"
	"github.com/labstack/echo/v4"
)

type handler struct {
	service TaskService
}

func New(s TaskService) h.Handler {
	return decorator.NewBindAndValidateDecorator(&handler{service: s})
}

type Request struct{}

func (h *handler) Handle(c echo.Context, in Request) error {
	tasks, err := h.service.ListTasks()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tasks)
}
