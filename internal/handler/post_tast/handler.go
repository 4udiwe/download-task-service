package post_task

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

type PostTaskRequest struct {
	URLs []string `json:"URLs" validate:"required"`
}

func (h *handler) Handle(c echo.Context, in PostTaskRequest) error {
	task, err := h.service.CreateTask(in.URLs)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, task)
}
