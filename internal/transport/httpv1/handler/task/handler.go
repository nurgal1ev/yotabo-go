package task

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/handler/common"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nurgal1ev/yotabo-go/internal/infrastructure/postgres"
	"github.com/nurgal1ev/yotabo-go/internal/models"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/middleware"
	"gorm.io/gorm"
)

func CreateTaskHandler(ctx context.Context, input *CreateTaskInput) (*common.HumaAPIResponse[TaskDTO], error) {
	userID := middleware.GetUserID(ctx)
	slog.Info("create task user", "userID", userID)

	var dueTime *time.Time
	if input.Body.DueDate != nil && *input.Body.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *input.Body.DueDate)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid date format, expected YYYY-MM-DD")
		}
		dueTime = &parsed
	}

	task := &models.Task{
		Name:        input.Body.Name,
		Description: input.Body.Description,
		Status:      input.Body.Status,
		Priority:    input.Body.Priority,
		DueDate:     dueTime,
		BoardID:     input.Body.BoardID,
		CreatedByID: uint(userID),
		UpdatedByID: uint(userID),
	}

	err := gorm.G[models.Task](postgres.Db).Create(ctx, task)
	if err != nil {
		slog.Error("failed create task", slog.String("error", err.Error()))
		return nil, err
	}

	var dueDateStr *string
	if task.DueDate != nil {
		s := task.DueDate.Format("2006-01-02")
		dueDateStr = &s
	}

	return common.NewHumaResponse(TaskDTO{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		BoardID:     task.BoardID,
		DueDate:     dueDateStr,
	}), nil
}

func GetTaskHandler(ctx context.Context, input *GetTaskInput) (*common.HumaAPIResponse[TaskDTO], error) {
	task, err := gorm.G[models.Task](postgres.Db).Preload("Subtasks", nil).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("task not found")
		}
		slog.Error("failed get task", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError("internal server error")
	}

	subtasks, err := gorm.G[models.Subtask](postgres.Db).Where("task_id = ?", task.ID).Find(ctx)
	if err != nil {
		slog.Error("failed get subtasks", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError("internal server error")
	}

	subtasksDTOs := make([]SubtasksDTO, len(subtasks))
	for i, subtask := range subtasks {
		subtasksDTOs[i] = SubtasksDTO{
			ID:        subtask.ID,
			Name:      subtask.Name,
			Completed: subtask.Completed,
		}
	}

	comments, err := gorm.G[models.Comment](postgres.Db).Where("task_id = ?", task.ID).Find(ctx)
	if err != nil {
		slog.Error("failed get comments", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError("internal server error")
	}

	commentsDTOs := make([]CommentsDTO, len(comments))
	for i, comment := range comments {
		commentsDTOs[i] = CommentsDTO{
			ID:      comment.ID,
			Message: comment.Message,
		}
	}

	var dueDateStr *string
	if task.DueDate != nil {
		s := task.DueDate.Format("2006-01-02")
		dueDateStr = &s
	}

	return common.NewHumaResponse(TaskDTO{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		DueDate:     dueDateStr,
		Subtasks:    subtasksDTOs,
		Comments:    commentsDTOs,
		BoardID:     task.BoardID,
	}), nil
}

func GetAllTasksHandler(ctx context.Context, input *GetAllTasksInput) (*common.HumaAPIResponse[[]TaskDTO], error) {
	userID := middleware.GetUserID(ctx)
	query := gorm.G[models.Task](postgres.Db).Preload("Subtasks", nil).Where("created_by_id = ?", userID)

	if input.BoardID.IsSet {
		query = query.Where("board_id = ?", input.BoardID.Value)
	}

	tasks, err := query.Find(ctx)
	if err != nil {
		slog.Error("failed get all tasks", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError("internal server error")
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
			BoardID:     task.BoardID,
		}
	}

	return common.NewHumaResponse(taskDTOs), nil
}

func UpdateTaskHandler(ctx context.Context, input *UpdateTaskInput) (*common.HumaAPIResponse[TaskDTO], error) {
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
		task.Description = input.Body.Description
	}
	if input.Body.Status != nil {
		task.Status = *input.Body.Status
	}
	if input.Body.Priority != nil {
		task.Priority = *input.Body.Priority
	}

	if input.Body.DueDate != nil {
		if *input.Body.DueDate == "" {
			task.DueDate = nil
		} else {
			parsed, err := time.Parse("2006-01-02", *input.Body.DueDate)
			if err != nil {
				return nil, huma.Error400BadRequest("invalid date format, expected YYYY-MM-DD")
			}
			task.DueDate = &parsed
		}
	}

	_, err = gorm.G[models.Task](postgres.Db).
		Where("id = ?", input.ID).
		Updates(ctx, task)

	if err != nil {
		return nil, huma.Error500InternalServerError("internal server error")
	}

	var dueDateStr *string
	if task.DueDate != nil {
		s := task.DueDate.Format("2006-01-02")
		dueDateStr = &s
	}

	return common.NewHumaResponse(TaskDTO{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		DueDate:     dueDateStr,
		BoardID:     task.BoardID,
	}), nil
}

func DeleteTaskHandler(ctx context.Context, input *DeleteTaskInput) (*common.HumaAPIResponse[any], error) {
	userID := middleware.GetUserID(ctx)

	taskID := input.ID

	result := postgres.Db.Delete(&models.Task{}, "id = ? AND created_by_id = ?", taskID, userID)

	if result.Error != nil {
		return nil, result.Error
	}

	return common.NewHumaResponse[any](nil, "task deleted successfully"), nil
}
