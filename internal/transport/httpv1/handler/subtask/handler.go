package subtask

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nurgal1ev/yotabo-go/internal/infrastructure/postgres"
	"github.com/nurgal1ev/yotabo-go/internal/models"
	"gorm.io/gorm"
)

func CreateSubtaskHandler(ctx context.Context, input *SubtaskResponse) (*CreateSubtaskOutput, error) {
	err := gorm.G[models.Subtask](postgres.Db).Create(ctx, &models.Subtask{
		Name:      input.Body.Name,
		Completed: false,
		TaskID:    input.Body.TaskID,
	})

	if err != nil {
		slog.Error("failed create subtask", slog.String("error", err.Error()))
		return nil, err
	}

	return &CreateSubtaskOutput{
		Status: http.StatusCreated,
		Body: struct {
			Message string `json:"message"`
		}{Message: "success"},
	}, nil
}

func UpdateSubtaskHandler(ctx context.Context, input *UpdateSubtaskInput) (*UpdateSubtaskOutput, error) {
	subtask, err := gorm.G[models.Subtask](postgres.Db).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("subtask not found")
		}
		return nil, err
	}

	if input.Body.Name != nil {
		subtask.Name = *input.Body.Name
	}

	if input.Body.Completed != nil {
		subtask.Completed = *input.Body.Completed
	}

	_, err = gorm.G[models.Subtask](postgres.Db).
		Where("id = ?", input.ID).
		Updates(ctx, subtask)

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &UpdateSubtaskOutput{
		Status: http.StatusOK,
		Body: SubtaskResponse{
			Body: SubtaskDTO{
				Name:      subtask.Name,
				Completed: subtask.Completed,
			},
		},
	}, nil
}
