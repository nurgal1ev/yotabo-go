package subtask

import (
	"context"
	"errors"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nurgal1ev/yotabo-go/internal/infrastructure/postgres"
	"github.com/nurgal1ev/yotabo-go/internal/models"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/handler/common"
	"gorm.io/gorm"
)

func CreateSubtaskHandler(ctx context.Context, input *CreateSubtaskInput) (*common.HumaAPIResponse[SubtaskDTO], error) {
	err := gorm.G[models.Subtask](postgres.Db).Create(ctx, &models.Subtask{
		Name:      input.Body.Name,
		Completed: false,
		TaskID:    input.TaskID,
	})

	if err != nil {
		slog.Error("failed create subtask", slog.String("error", err.Error()))
		return nil, err
	}

	return common.NewHumaResponse(SubtaskDTO{
		Name:      input.Body.Name,
		Completed: input.Body.Completed,
		TaskID:    input.TaskID,
	}), nil
}

func UpdateSubtaskHandler(ctx context.Context, input *UpdateSubtaskInput) (*common.HumaAPIResponse[SubtaskDTO], error) {
	subtask, err := gorm.G[models.Subtask](postgres.Db).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("subtask not found")
		}
		return nil, err
	}

	subtask.Completed = *input.Body.Completed

	_, err = gorm.G[models.Subtask](postgres.Db).
		Where("id = ?", input.ID).
		Updates(ctx, subtask)

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return common.NewHumaResponse(SubtaskDTO{
		Name:      subtask.Name,
		Completed: subtask.Completed,
		TaskID:    subtask.TaskID,
	}), nil
}
