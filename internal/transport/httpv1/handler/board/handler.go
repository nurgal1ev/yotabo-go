package board

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

func CreateBoardHandler(ctx context.Context, input *CreateBoardInput) (*common.HumaAPIResponse[BoardDTO], error) {
	userID := middleware.GetUserID(ctx)

	board := &models.Board{
		Name:        input.Body.Name,
		CreatedByID: uint(userID),
		FolderID:    input.Body.FolderID,
	}

	err := gorm.G[models.Board](postgres.Db).Create(ctx, board)
	if err != nil {
		slog.Error("failed create board", slog.String("error", err.Error()))
		return nil, err
	}

	return common.NewHumaResponse(BoardDTO{
		ID:       board.ID,
		Name:     board.Name,
		FolderID: board.FolderID,
	}), nil
}

func GetBoardHandler(ctx context.Context, input *GetBoardInput) (*common.HumaAPIResponse[BoardDTO], error) {
	board, err := gorm.G[models.Board](postgres.Db).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("board not found")
		}
		slog.Error("failed get board", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError("internal server error")
	}

	return common.NewHumaResponse(BoardDTO{
		ID:   input.ID,
		Name: board.Name,
	}), nil
}

func GetAllBoardsHandler(ctx context.Context, input *GetAllBoardsInput) (*common.HumaAPIResponse[[]BoardDTO], error) {
	userID := middleware.GetUserID(ctx)

	boards, err := gorm.G[models.Board](postgres.Db).Where("created_by_id = ?", userID).Find(ctx)
	if err != nil {
		slog.Error("failed get all boards", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError("internal server error")
	}

	boardDTOs := make([]BoardDTO, len(boards))
	for i, board := range boards {
		boardDTOs[i] = BoardDTO{
			ID:       board.ID,
			Name:     board.Name,
			FolderID: board.FolderID,
		}
	}

	return common.NewHumaResponse(boardDTOs), nil
}

func DeleteBoardHandler(ctx context.Context, input *DeleteBoardInput) (*common.HumaAPIResponse[any], error) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		return nil, huma.Error401Unauthorized("unauthorized")
	}

	boardID := input.ID

	result := postgres.Db.Delete(&models.Board{}, "id = ? AND created_by_id = ?", boardID, userID)

	if result.Error != nil {
		return nil, result.Error
	}

	return common.NewHumaResponse[any](nil, "board deleted successfully"), nil
}

func UpdateBoardHandler(ctx context.Context, input *UpdateBoardInput) (*common.HumaAPIResponse[BoardDTO], error) {
	board, err := gorm.G[models.Board](postgres.Db).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("board not found")
		}
		return nil, err
	}

	if input.Body.Name != nil {
		board.Name = *input.Body.Name
	}

	if input.Body.FolderID != nil {
		board.FolderID = input.Body.FolderID
	}

	_, err = gorm.G[models.Board](postgres.Db).
		Where("id = ?", input.ID).
		Updates(ctx, board)

	return common.NewHumaResponse(BoardDTO{
		ID:       board.ID,
		Name:     board.Name,
		FolderID: board.FolderID,
	}), nil
}
