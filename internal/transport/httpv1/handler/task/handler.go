package task

import (
	"context"
	"github.com/nurgal1ev/yotabo-go/internal/infrastructure/postgres"
	"github.com/nurgal1ev/yotabo-go/internal/models"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/middleware"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

type TaskResponse struct {
	Body struct {
		Name        string `json:"name" minLength:"1" maxLength:"55" pattern:"^[a-zA-Zа-яА-Я0-9\\s]+$"`
		Description string `json:"description" maxLength:"10000" pattern:"^[a-zA-Zа-яА-Я0-9\\s]+$"`
		Status      string `json:"status" enum:"backlog,in_progress,review,done"`
		Priority    string `json:"priority" enum:"easy,medium,hard"`
	}
}
type CreateTaskOutput struct {
	Status int
	Body   struct {
		Message string `json:"message"`
	}
}

type GetTaskInput struct {
	Params struct {
		ID uint `path:"id"` // {id} в URL
	}
}
type GetTaskOutput struct {
	Status int `status:"201"`
	Body   TaskResponse
}

func CreateTaskHandler(ctx context.Context, input *TaskResponse) (*CreateTaskOutput, error) {
	userID := middleware.GetUserID(ctx)

	err := gorm.G[models.Task](postgres.Db).Create(ctx, &models.Task{
		Name:        input.Body.Name,
		Description: input.Body.Description,
		Status:      input.Body.Status,
		Priority:    input.Body.Priority,
		CreatedByID: uint(userID),
		UpdatedByID: uint(userID),
	})

	if err != nil {
		slog.Error("failed create task", slog.String("error", err.Error()))
		return nil, err
	}

	return &CreateTaskOutput{
		Status: http.StatusCreated,
		Body: struct {
			Message string `json:"message"`
		}{Message: "success"}}, nil
}

/*func GetTaskHandler(ctx context.Context, input *GetTaskInput) (*GetTaskOutput, error) {
	task, err := tasks.GetTask(ctx, input.Params.ID)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	return &GetTaskOutput{
		Status: 200,
		Body: TaskResponse{
			Name:        task.Name,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
		},
	}, nil
}*/
