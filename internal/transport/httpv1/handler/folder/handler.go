package folder

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

func CreateFolderHandler(ctx context.Context, input *FolderResponse) (*CreateFolderOutput, error) {
	userID := middleware.GetUserID(ctx)

	err := gorm.G[models.Folder](postgres.Db).Create(ctx, &models.Folder{
		Name:        input.Body.Name,
		CreatedByID: uint(userID),
	})

	if err != nil {
		slog.Error("failed create folder", slog.String("error", err.Error()))
		return nil, err
	}

	return &CreateFolderOutput{
		Status: http.StatusCreated,
		Body: struct {
			Message string `json:"message"`
		}{Message: "success"},
	}, nil
}

func DeleteFolderHandler(ctx context.Context, input *DeleteFolderInput) (*DeleteFolderOutput, error) {
	userID := middleware.GetUserID(ctx)

	if userID == 0 {
		return nil, huma.Error401Unauthorized("unauthorized")
	}

	folderID := input.ID

	result := postgres.Db.Delete(&models.Folder{}, "id = ? AND created_by_id = ?", folderID, userID)

	if result.Error != nil {
		return nil, result.Error
	}

	return &DeleteFolderOutput{
		Status:  http.StatusOK,
		Message: "folder deleted successfully",
	}, nil
}

func GetFolderHandler(ctx context.Context, input *GetFolderInput) (*GetFolderOutput, error) {
	folder, err := gorm.G[models.Folder](postgres.Db).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("folder not found")
		}
		slog.Error("failed get folder", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &GetFolderOutput{
		Status: http.StatusOK,
		Body: FolderResponse{
			Body: FolderDTO{
				Name: folder.Name,
			},
		},
	}, nil
}

func GetAllFoldersHandler(ctx context.Context, input *GetAllFoldersInput) (*GetAllFoldersOutput, error) {
	userID := middleware.GetUserID(ctx)

	folders, err := gorm.G[models.Folder](postgres.Db).Where("created_by_id = ?", userID).Find(ctx)
	if err != nil {
		slog.Error("failed get all folders", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError(err.Error())
	}

	foldersDTOs := make([]FolderDTO, len(folders))
	for i, folder := range folders {
		foldersDTOs[i] = FolderDTO{
			Name: folder.Name,
		}
	}

	return &GetAllFoldersOutput{
		Status: http.StatusOK,
		Body: struct {
			Folders []FolderDTO `json:"folders"`
		}{
			Folders: foldersDTOs,
		},
	}, nil
}

func UpdateFolderHandler(ctx context.Context, input *UpdateFolderInput) (*UpdateFolderOutput, error) {
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

	return &UpdateFolderOutput{
		Status: http.StatusOK,
		Body: FolderResponse{
			Body: FolderDTO{
				Name: folder.Name,
			},
		},
	}, nil
}
