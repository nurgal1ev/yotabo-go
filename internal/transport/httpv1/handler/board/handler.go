package board

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nurgal1ev/yotabo-go/internal/infrastructure/postgres"
	"github.com/nurgal1ev/yotabo-go/internal/models"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/handler/common"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/middleware"
	"gorm.io/gorm"
)

// CRUD
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
	query := gorm.G[models.Board](postgres.Db).Where("created_by_id = ?", userID)

	if input.FolderID.IsSet {
		query = query.Where("folder_id = ?", input.FolderID.Value)
	}

	boards, err := query.Find(ctx)
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

// INVITE
func CreateInviteHandler(ctx context.Context, input *CreateInvitationInput) (*common.HumaAPIResponse[InvitationDTO], error) {
	userID := middleware.GetUserID(ctx)
	boardID := input.BoardID
	username := input.Body.Username

	board, err := gorm.G[models.Board](postgres.Db).Where("id = ?", boardID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("board not found")
		}
		return nil, err
	}

	invitedUser, err := gorm.G[models.User](postgres.Db).Where("username = ?", username).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("user not found")
		}
		return nil, err
	}

	if uint(userID) == invitedUser.ID {
		return nil, huma.Error400BadRequest("you cannot invite yourself")
	}

	_, err = gorm.G[models.BoardMember](postgres.Db).Where("board_id = ? AND user_id = ?", boardID, invitedUser.ID).First(ctx)

	if err == nil {
		return nil, huma.Error409Conflict("user is already a member of this board")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("failed check membership", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError("internal server error")
	}

	if uint(userID) != board.CreatedByID {
		return nil, huma.Error403Forbidden("only board owner can invite members")
	}

	var invite models.Invitation
	err = postgres.Db.Where("board_id = ? AND user_id = ? AND status = ?", boardID, invitedUser.ID, "pending").First(&invite).Error

	if err == nil {
		return nil, huma.Error409Conflict("invitation already sent to this user")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("failed check existing invitation", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError("internal server error")
	}

	var inviterUser models.User
	inviterName := ""
	err = postgres.Db.First(&inviterUser, userID).Error
	if err == nil {
		inviterName = inviterUser.Username
	} else {
		slog.Warn("failed to get inviter name", slog.String("error", err.Error()))
	}

	invitation := models.Invitation{
		BoardID:   board.ID,
		UserID:    invitedUser.ID,
		InviterID: uint(userID),
		Status:    "pending",
	}

	err = postgres.Db.Create(&invitation).Error
	if err != nil {
		slog.Error("failed create invitation", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError("internal server error")
	}

	return common.NewHumaResponse(InvitationDTO{
		ID:        invitation.ID,
		BoardID:   invitation.BoardID,
		BoardName: board.Name,
		InviterID: invitation.InviterID,
		Inviter:   inviterName,
		Status:    invitation.Status,
		CreatedAt: invitation.CreatedAt.Format(time.RFC3339),
	}), nil
}
