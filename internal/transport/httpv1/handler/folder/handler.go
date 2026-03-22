package folder

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

func CreateFolderHandler(ctx context.Context, input *CreateFolderInput) (*common.HumaAPIResponse[FolderDTO], error) {
	userID := middleware.GetUserID(ctx)

	folder := &models.Folder{
		Name:        input.Body.Name,
		CreatedByID: uint(userID),
	}
	err := gorm.G[models.Folder](postgres.Db).Create(ctx, folder)

	if err != nil {
		slog.Error("failed create folder", slog.String("error", err.Error()))
		return nil, err
	}

	return common.NewHumaResponse(FolderDTO{
		ID:   folder.ID,
		Name: folder.Name,
	}), nil
}

func DeleteFolderHandler(ctx context.Context, input *DeleteFolderInput) (*common.HumaAPIResponse[any], error) {
	userID := middleware.GetUserID(ctx)

	if userID == 0 {
		return nil, huma.Error401Unauthorized("unauthorized")
	}

	folderID := input.ID

	result := postgres.Db.Delete(&models.Folder{}, "id = ? AND created_by_id = ?", folderID, userID)

	if result.Error != nil {
		return nil, result.Error
	}

	return common.NewHumaResponse[any](nil, "folder deleted successfully"), nil
}

func GetFolderHandler(ctx context.Context, input *GetFolderInput) (*common.HumaAPIResponse[FolderDTO], error) {
	folder, err := gorm.G[models.Folder](postgres.Db).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("folder not found")
		}
		slog.Error("failed get folder", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError(err.Error())
	}

	boards, err := gorm.G[models.Board](postgres.Db).Where("folder_id = ?", folder.ID).Find(ctx)
	if err != nil {
		slog.Error("failed get all boards", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError(err.Error())
	}

	boardDTOs := make([]FolderBoardDTO, len(boards))
	for i, board := range boards {
		boardDTOs[i] = FolderBoardDTO{
			Name: board.Name,
		}
	}

	return common.NewHumaResponse(FolderDTO{
		ID:     folder.ID,
		Name:   folder.Name,
		Boards: boardDTOs,
	}), nil
}

func GetAllFoldersHandler(ctx context.Context, input *GetAllFoldersInput) (*common.HumaAPIResponse[[]FolderDTO], error) {
	userID := middleware.GetUserID(ctx)

	folders, err := gorm.G[models.Folder](postgres.Db).Preload("Boards", nil).Where("created_by_id = ?", userID).Find(ctx)
	if err != nil {
		slog.Error("failed get all folders", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError(err.Error())
	}

	foldersDTOs := make([]FolderDTO, len(folders))
	for i, folder := range folders {
		boardDTOs := make([]FolderBoardDTO, 0, len(folder.Boards))
		for _, board := range folder.Boards {
			boardDTOs = append(boardDTOs, FolderBoardDTO{
				ID:   board.ID,
				Name: board.Name,
			})
		}
		foldersDTOs[i] = FolderDTO{
			ID:     folder.ID,
			Name:   folder.Name,
			Boards: boardDTOs,
		}
	}

	return common.NewHumaResponse(foldersDTOs), nil
}

func UpdateFolderHandler(ctx context.Context, input *UpdateFolderInput) (*common.HumaAPIResponse[FolderDTO], error) {
	userID := middleware.GetUserID(ctx)
	folder, err := gorm.G[models.Folder](postgres.Db).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("folder not found")
		}
		return nil, err
	}

	if userID == 0 {
		return nil, huma.Error401Unauthorized("unauthorized")
	}

	if input.Body.Name != nil {
		folder.Name = *input.Body.Name
	}

	_, err = gorm.G[models.Folder](postgres.Db).
		Where("id = ?", input.ID).
		Updates(ctx, folder)

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return common.NewHumaResponse(FolderDTO{
		Name: folder.Name,
	}), nil
}
