package comment

import (
	"context"
	"log/slog"

	"github.com/nurgal1ev/yotabo-go/internal/infrastructure/postgres"
	"github.com/nurgal1ev/yotabo-go/internal/models"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/handler/common"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/middleware"
	"gorm.io/gorm"
)

func CreateCommentHandler(ctx context.Context, input *CreateCommentInput) (*common.HumaAPIResponse[CommentDTO], error) {
	userID := middleware.GetUserID(ctx)
	err := gorm.G[models.Comment](postgres.Db).Create(ctx, &models.Comment{
		Message:  input.Body.Message,
		AuthorID: uint(userID),
		TaskID:   input.TaskID,
	})

	if err != nil {
		slog.Error("failed create comment", slog.String("error", err.Error()))
		return nil, err
	}

	return common.NewHumaResponse(CommentDTO{
		Message:  input.Body.Message,
		AuthorID: uint(userID),
		TaskID:   input.TaskID,
	}), nil
}
