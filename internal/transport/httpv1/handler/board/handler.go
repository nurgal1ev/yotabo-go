package board

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

func CreateBoardHandler(ctx context.Context, input *CreateBoardInput) (*CreateBoardOutput, error) {
	userID := middleware.GetUserID(ctx)

	err := gorm.G[models.Board](postgres.Db).Create(ctx, &models.Board{
		Name:        input.Body.Name,
		CreatedByID: uint(userID),
		FolderID:    input.Body.FolderID,
	})

	if err != nil {
		slog.Error("failed create board", slog.String("error", err.Error()))
		return nil, err
	}

	return &CreateBoardOutput{
		Status: http.StatusCreated,
		Body: struct {
			Message string `json:"message"`
		}{Message: "success"},
	}, nil
}

func GetBoardHandler(ctx context.Context, input *GetBoardInput) (*GetBoardOutput, error) {
	board, err := gorm.G[models.Board](postgres.Db).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("board not found")
		}
		slog.Error("failed get board", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &GetBoardOutput{
		Status: http.StatusOK,
		Body: BoardDTO{
			ID:   input.ID,
			Name: board.Name,
		},
	}, nil
}

func GetAllBoardsHandler(ctx context.Context, input *GetAllBoardsInput) (*GetAllBoardsOutput, error) {
	userID := middleware.GetUserID(ctx)

	boards, err := gorm.G[models.Board](postgres.Db).Where("created_by_id = ?", userID).Find(ctx)
	if err != nil {
		slog.Error("failed get all boards", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError(err.Error())
	}

	boardDTOs := make([]BoardDTO, len(boards))
	for i, board := range boards {
		boardDTOs[i] = BoardDTO{
			ID:       board.ID,
			Name:     board.Name,
			FolderID: board.FolderID,
		}
	}

	return &GetAllBoardsOutput{
		Status: http.StatusOK,
		Body: struct {
			Boards []BoardDTO `json:"boards"`
		}{
			Boards: boardDTOs,
		},
	}, nil
}

func DeleteBoardHandler(ctx context.Context, input *DeleteBoardInput) (*DeleteBoardOutput, error) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		return nil, huma.Error401Unauthorized("unauthorized")
	}

	boardID := input.ID

	result := postgres.Db.Delete(&models.Board{}, "id = ? AND created_by_id = ?", boardID, userID)

	if result.Error != nil {
		return nil, result.Error
	}

	return &DeleteBoardOutput{
		Status:  http.StatusOK,
		Message: "board deleted successfully",
	}, nil
}

func UpdateBoardHandler(ctx context.Context, input *UpdateBoardInput) (*UpdateBoardOutput, error) {
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

	return &UpdateBoardOutput{
		Status: http.StatusOK,
		Body: BoardDTO{
			ID:       board.ID,
			Name:     board.Name,
			FolderID: board.FolderID,
		},
	}, nil
}
