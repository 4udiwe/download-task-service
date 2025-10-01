package get_task_by_id

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

type GetTaskByIDRequest struct {
	ID string `path:"taskID" validate:"required"`
}

func (h *handler) Handle(c echo.Context, in GetTaskByIDRequest) error {
	task, err := h.service.GetTask(in.ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, task)
}
