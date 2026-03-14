package task

import (
	"context"
	"errors"
	"github.com/danielgtaylor/huma/v2"
	"github.com/nurgal1ev/yotabo-go/internal/infrastructure/postgres"
	"github.com/nurgal1ev/yotabo-go/internal/models"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/middleware"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

// TODO: прокинуть айди
func CreateTaskHandler(ctx context.Context, input *TaskResponse) (*CreateTaskOutput, error) {
	userID := middleware.GetUserID(ctx)
	slog.Info("create task user", "userID", userID)

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

func GetTaskHandler(ctx context.Context, input *GetTaskInput) (*GetTaskOutput, error) {
	task, err := gorm.G[models.Task](postgres.Db).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("task not found")
		}
		slog.Error("failed get task", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &GetTaskOutput{
		Status: http.StatusOK,
		Body: TaskDTO{
			ID:          task.ID,
			Name:        task.Name,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
		},
	}, nil
}

func GetAllTasksHandler(ctx context.Context, input *GetAllTasksInput) (*GetAllTasksOutput, error) {
	userID := middleware.GetUserID(ctx)

	if userID == 0 {
		return nil, huma.Error401Unauthorized("unauthorized")
	}

	tasks, err := gorm.G[models.Task](postgres.Db).Where("created_by_id = ?", userID).Find(ctx)
	if err != nil {
		slog.Error("failed get all tasks", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError(err.Error())
	}

	taskDTOs := make([]TaskDTO, len(tasks))
	for i, task := range tasks {
		taskDTOs[i] = TaskDTO{
			ID:          task.ID,
			Name:        task.Name,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
		}
	}

	return &GetAllTasksOutput{
		Status: http.StatusOK,
		Body: struct {
			Tasks []TaskDTO `json:"tasks"`
		}{
			Tasks: taskDTOs,
		},
	}, nil
}

func UpdateTaskHandler(ctx context.Context, input *UpdateTaskInput) (*UpdateTaskOutput, error) {
	task, err := gorm.G[models.Task](postgres.Db).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("task not found")
		}
		return nil, err
	}

	if input.Body.Name != nil {
		task.Name = *input.Body.Name
	}
	if input.Body.Description != nil {
		task.Description = *input.Body.Description
	}
	if input.Body.Status != nil {
		task.Status = *input.Body.Status
	}
	if input.Body.Priority != nil {
		task.Priority = *input.Body.Priority
	}

	_, err = gorm.G[models.Task](postgres.Db).
		Where("id = ?", input.ID).
		Updates(ctx, task)

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &UpdateTaskOutput{
		Status: http.StatusOK,
		Body: TaskResponse{
			Body: TaskDTO{
				Name:        task.Name,
				Description: task.Description,
				Status:      task.Status,
				Priority:    task.Priority,
			},
		},
	}, nil
}

func DeleteTaskHandler(ctx context.Context, input *DeleteTaskInput) (*DeleteTaskOutput, error) {
	userID := middleware.GetUserID(ctx)

	if userID == 0 {
		return nil, huma.Error401Unauthorized("unauthorized")
	}

	taskID := input.ID

	result := postgres.Db.Delete(&models.Task{}, "id = ? AND created_by_id = ?", taskID, userID)

	if result.Error != nil {
		return nil, result.Error
	}

	return &DeleteTaskOutput{
		Status:  http.StatusOK,
		Message: "task deleted successfully",
	}, nil
}
