package note

import (
	"context"
	"errors"
	"fmt"

	transaction "github.com/khaingminhtun/relio-backend/internal/shared/dbutils"
	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
	"gorm.io/gorm"
)

// ============================================================
// Interface
// ============================================================

type Repository interface {
	Create(
		ctx context.Context,
		note *Note,
	) error

	GetByID(
		ctx context.Context,
		userID int64,
		noteID int64,
	) (*Note, error)

	List(
		ctx context.Context,
		userID int64,
		filter ListNotesFilter,
	) ([]Note, error)

	Update(
		ctx context.Context,
		note *Note,
	) error

	Delete(
		ctx context.Context,
		userID int64,
		noteID int64,
	) error
}

// ============================================================
// Implementation
// ============================================================

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(
	ctx context.Context,
	note *Note,
) error {
	return transaction.DB(ctx, r.db).
		Create(note).
		Error
}

func (r *repository) GetByID(
	ctx context.Context,
	userID int64,
	noteID int64,
) (*Note, error) {
	var note Note

	err := transaction.DB(ctx, r.db).
		Where("id = ?", noteID).
		Where("user_id = ?", userID).
		First(&note).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeNoteNotFound,
				"note not found",
				err,
			)
		}

		return nil, err
	}

	return &note, nil
}

func (r *repository) List(
	ctx context.Context,
	userID int64,
	filter ListNotesFilter,
) ([]Note, error) {
	var notes []Note

	q := transaction.DB(ctx, r.db).
		Where("user_id = ?", userID)

	// --------------------------------------------------------
	// Optional filters — always scoped to the owning user
	// --------------------------------------------------------

	if filter.Mood != nil {
		q = q.Where("mood = ?", *filter.Mood)
	}

	if filter.IsPinned != nil {
		q = q.Where("is_pinned = ?", *filter.IsPinned)
	}

	if filter.IsArchived != nil {
		q = q.Where("is_archived = ?", *filter.IsArchived)
	}

	if filter.Search != nil && *filter.Search != "" {
		pattern := fmt.Sprintf("%%%s%%", *filter.Search)
		q = q.Where("title ILIKE ? OR content ILIKE ?", pattern, pattern)
	}

	err := q.
		Order("created_at DESC").
		Find(&notes).
		Error

	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *repository) Update(
	ctx context.Context,
	note *Note,
) error {
	return transaction.DB(ctx, r.db).
		Save(note).
		Error
}

func (r *repository) Delete(
	ctx context.Context,
	userID int64,
	noteID int64,
) error {
	result := transaction.DB(ctx, r.db).
		Where("id = ?", noteID).
		Where("user_id = ?", userID).
		Delete(&Note{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperror.New(
			apperror.CodeNoteNotFound,
			"note not found",
			gorm.ErrRecordNotFound,
		)
	}

	return nil
}
