package subtask

import (
	"context"
	"errors"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nurgal1ev/yotabo-go/internal/infrastructure/postgres"
	"github.com/nurgal1ev/yotabo-go/internal/models"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/handler/common"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/middleware"
	"gorm.io/gorm"
)

func CreateSubtaskHandler(ctx context.Context, input *CreateSubtaskInput) (*common.HumaAPIResponse[SubtaskDTO], error) {
	subtask := &models.Subtask{
		Name:      input.Body.Name,
		Completed: false,
		TaskID:    input.TaskID,
	}

	err := gorm.G[models.Subtask](postgres.Db).Create(ctx, subtask)

	if err != nil {
		slog.Error("failed create subtask", slog.String("error", err.Error()))
		return nil, err
	}

	return common.NewHumaResponse(SubtaskDTO{
		ID:        subtask.ID,
		Name:      subtask.Name,
		Completed: subtask.Completed,
		TaskID:    subtask.TaskID,
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

	return common.NewHumaResponse(SubtaskDTO{
		Name:      subtask.Name,
		Completed: subtask.Completed,
		TaskID:    subtask.TaskID,
	}), nil
}

func DeleteSubtaskHandler(ctx context.Context, input *DeleteSubtaskInput) (*common.HumaAPIResponse[any], error) {
	userID := middleware.GetUserID(ctx)

	result := postgres.Db.Delete(&models.Subtask{}, "id = ? AND task_id IN (SELECT id FROM tasks WHERE created_by_id = ?)", input.ID, userID)

	if result.Error != nil {
		slog.Error("failed delete subtask", slog.String("error", result.Error.Error()))
		return nil, huma.Error500InternalServerError("failed to delete subtask")
	}

	if result.RowsAffected == 0 {
		return nil, huma.Error404NotFound("subtask not found")
	}

	return common.NewHumaResponse[any](nil, "subtask deleted successfully"), nil
}
