package task

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nurgal1ev/yotabo-go/internal/infrastructure/postgres"
	"github.com/nurgal1ev/yotabo-go/internal/models"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/middleware"
	"gorm.io/gorm"
)

// TODO: прокинуть айди
func CreateTaskHandler(ctx context.Context, input *CreateTaskInput) (*CreateTaskOutput, error) {
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
	task, err := gorm.G[models.Task](postgres.Db).Preload("Subtasks", nil).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("task not found")
		}
		slog.Error("failed get task", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError(err.Error())
	}

	subtasks, err := gorm.G[models.Subtask](postgres.Db).Where("task_id = ?", task.ID).Find(ctx)
	if err != nil {
		slog.Error("failed get subtasks", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError(err.Error())
	}

	subtasksDTOs := make([]SubtasksDTO, len(subtasks))
	for i, subtask := range subtasks {
		subtasksDTOs[i] = SubtasksDTO{
			ID:        subtask.ID,
			Name:      subtask.Name,
			Completed: subtask.Completed,
		}
	}

	return &GetTaskOutput{
		Status: http.StatusOK,
		Body: TaskDTO{
			ID:          task.ID,
			Name:        task.Name,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
			Subtasks:    subtasksDTOs,
		},
	}, nil
}

func GetAllTasksHandler(ctx context.Context, input *GetAllTasksInput) (*GetAllTasksOutput, error) {
	userID := middleware.GetUserID(ctx)

	if userID == 0 {
		return nil, huma.Error401Unauthorized("unauthorized")
	}

	tasks, err := gorm.G[models.Task](postgres.Db).Preload("Subtasks", nil).Where("created_by_id = ?", userID).Find(ctx)
	if err != nil {
		slog.Error("failed get all tasks", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError(err.Error())
	}

	taskDTOs := make([]TaskDTO, len(tasks))
	for i, task := range tasks {
		subtaskDTOs := make([]SubtasksDTO, 0, len(task.Subtasks))
		for _, subtask := range task.Subtasks {
			subtaskDTOs = append(subtaskDTOs, SubtasksDTO{
				ID:        subtask.ID,
				Name:      subtask.Name,
				Completed: subtask.Completed,
			})
		}
		taskDTOs[i] = TaskDTO{
			ID:          task.ID,
			Name:        task.Name,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
			Subtasks:    subtaskDTOs,
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
