package tasks

import (
	"context"
	"errors"
	"github.com/nurgal1ev/yotabo-go/internal/infrastructure/postgres"
	"github.com/nurgal1ev/yotabo-go/internal/models"
	"gorm.io/gorm"
)

type TaskData struct {
	UserID      uint
	Name        string
	Description string
	Status      string
	Priority    string
}

func CreateTask(ctx context.Context, task *TaskData) error {
	if task.Name == "" {
		return errors.New("task name is required")
	}

	err := gorm.G[models.Task](postgres.Db).Create(ctx, &models.Task{
		Name:        task.Name,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
	})

	if err != nil {
		return err
	}
	return nil
}

func DeleteTask(ctx context.Context, id uint) error {
	_, err := gorm.G[models.Task](postgres.Db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTask(ctx context.Context, id uint, task *TaskData) error {
	updates := models.Task{
		Name:        task.Name,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
	}

	_, err := gorm.G[models.Task](postgres.Db).Where("id = ?", id).Updates(ctx, updates)

	if err != nil {
		return err
	}

	return nil
}

func GetTask(ctx context.Context, id uint) (*TaskData, error) {
	task, err := gorm.G[models.Task](postgres.Db).Where("id =?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &TaskData{
		Name:        task.Name,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
	}, nil
}

func GetAllTasks(ctx context.Context) ([]TaskData, error) {
	var tasks []TaskData
	result := postgres.Db.WithContext(ctx).Find(&tasks)

	if result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil
}
