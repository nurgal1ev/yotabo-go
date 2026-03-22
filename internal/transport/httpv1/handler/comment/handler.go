package comment

import "C"
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

func CreateCommentHandler(ctx context.Context, input *CreateCommentInput) (*common.HumaAPIResponse[CommentDTO], error) {
	userID := middleware.GetUserID(ctx)

	comment := &models.Comment{
		Message:  input.Body.Message,
		AuthorID: uint(userID),
		TaskID:   input.TaskID,
	}

	err := gorm.G[models.Comment](postgres.Db).Create(ctx, comment)

	if err != nil {
		slog.Error("failed create comment", slog.String("error", err.Error()))
		return nil, err
	}

	return common.NewHumaResponse(CommentDTO{
		ID:       comment.ID,
		Message:  comment.Message,
		AuthorID: comment.AuthorID,
		TaskID:   comment.TaskID,
	}), nil
}

func UpdateCommentHandler(ctx context.Context, input *UpdateCommentInput) (*common.HumaAPIResponse[CommentDTO], error) {
	comment, err := gorm.G[models.Comment](postgres.Db).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("comment not found")
		}
		return nil, err
	}

	if input.Body.Message != nil {
		comment.Message = *input.Body.Message
	}

	_, err = gorm.G[models.Comment](postgres.Db).Where("id = ?", input.ID).Updates(ctx, comment)
	if err != nil {
		slog.Error("failed update comment", slog.String("error", err.Error()))
	}

	return common.NewHumaResponse(CommentDTO{
		Message: comment.Message,
	}), nil
}

func DeleteCommentHandler(ctx context.Context, input *DeleteCommentInput) (*common.HumaAPIResponse[any], error) {
	userID := middleware.GetUserID(ctx)

	result := postgres.Db.
		Table("comments").
		Where("comments.id = ? AND comments.task_id IN (SELECT id FROM tasks WHERE created_by_id = ?)", input.ID, userID).
		Delete(&models.Comment{})

	if result.Error != nil {
		slog.Error("failed delete comment", slog.String("error", result.Error.Error()))
		return nil, huma.Error500InternalServerError("failed to delete comment")
	}

	if result.RowsAffected == 0 {
		return nil, huma.Error404NotFound("comment not found")
	}

	return common.NewHumaResponse[any](nil, "comment deleted successfully"), nil
}

func GetCommentHandler(ctx context.Context, input *GetCommentInput) (*common.HumaAPIResponse[CommentDTO], error) {
	comment, err := gorm.G[models.Comment](postgres.Db).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("comment not found")
		}
	}

	return common.NewHumaResponse(CommentDTO{
		ID:       comment.ID,
		Message:  comment.Message,
		TaskID:   comment.TaskID,
		AuthorID: comment.AuthorID,
	}), nil
}

func GetCommentsByTaskHandler(ctx context.Context, input *GetCommentsByTaskInput) (*common.HumaAPIResponse[[]CommentDTO], error) {
	comments, err := gorm.G[models.Comment](postgres.Db).Where("task_id = ?", input.TaskID).Find(ctx)
	if err != nil {
		slog.Error("failed get comments", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError(err.Error())
	}

	commentsDTOs := make([]CommentDTO, len(comments))
	for i, comment := range comments {
		commentsDTOs[i] = CommentDTO{
			ID:       comment.ID,
			Message:  comment.Message,
			TaskID:   comment.TaskID,
			AuthorID: comment.AuthorID,
		}
	}

	return common.NewHumaResponse(commentsDTOs), nil
}
