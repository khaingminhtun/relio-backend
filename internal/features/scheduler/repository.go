package scheduler

import (
	"context"
	"errors"
	"time"

	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
	"gorm.io/gorm"
)

type ScheduledJobRepository interface {
	Create(
		ctx context.Context,
		job *ScheduledJob,
	) error

	GetByID(
		ctx context.Context,
		id int64,
	) (*ScheduledJob, error)

	ListDueJobs(
		ctx context.Context,
		now time.Time,
		limit int,
	) ([]ScheduledJob, error)

	Claim(
		ctx context.Context,
		id int64,
		startedAt time.Time,
	) error

	MarkCompleted(
		ctx context.Context,
		id int64,
		completedAt time.Time,
	) error

	MarkFailed(
		ctx context.Context,
		id int64,
		errMessage string,
	) error

	Cancel(
		ctx context.Context,
		id int64,
	) error

	ResetStaleJobs(
		ctx context.Context,
		before time.Time,
	) error
}

type scheduledJobRepository struct {
	db *gorm.DB
}

func NewScheduledJobRepository(
	db *gorm.DB,
) ScheduledJobRepository {
	return &scheduledJobRepository{
		db: db,
	}
}

func (r *scheduledJobRepository) Create(
	ctx context.Context,
	job *ScheduledJob,
) error {
	return r.db.WithContext(ctx).
		Create(job).
		Error
}

func (r *scheduledJobRepository) GetByID(
	ctx context.Context,
	id int64,
) (*ScheduledJob, error) {
	var job ScheduledJob

	err := r.db.WithContext(ctx).
		First(&job, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeJobNotFound,
				"sheduled job not found",
				nil,
			)
		}

		return nil, err
	}

	return &job, nil
}

func (r *scheduledJobRepository) ListDueJobs(
	ctx context.Context,
	now time.Time,
	limit int,
) ([]ScheduledJob, error) {
	var jobs []ScheduledJob

	if limit <= 0 {
		limit = 100
	}

	err := r.db.WithContext(ctx).
		Where("status = ?", JobStatusPending).
		Where("scheduled_at <= ?", now).
		Order("scheduled_at ASC").
		Limit(limit).
		Find(&jobs).
		Error

	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *scheduledJobRepository) Claim(
	ctx context.Context,
	id int64,
	startedAt time.Time,
) error {
	result := r.db.WithContext(ctx).
		Model(&ScheduledJob{}).
		Where("id = ?", id).
		Where("status = ?", JobStatusPending).
		Updates(map[string]interface{}{
			"status":     JobStatusProcessing,
			"started_at": startedAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperror.New(
			apperror.CodeJobNotClaimed,
			"scheduled job could not be claimed",
			nil,
		)
	}

	return nil
}

func (r *scheduledJobRepository) MarkCompleted(
	ctx context.Context,
	id int64,
	completedAt time.Time,
) error {
	result := r.db.WithContext(ctx).
		Model(&ScheduledJob{}).
		Where("id = ?", id).
		Where("status = ?", JobStatusProcessing).
		Updates(map[string]interface{}{
			"status":       JobStatusCompleted,
			"completed_at": completedAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *scheduledJobRepository) MarkFailed(
	ctx context.Context,
	id int64,
	errMessage string,
) error {
	result := r.db.WithContext(ctx).
		Model(&ScheduledJob{}).
		Where("id = ?", id).
		Where("status = ?", JobStatusProcessing).
		Updates(map[string]interface{}{
			"status":     JobStatusFailed,
			"attempts":   gorm.Expr("attempts + ?", 1),
			"last_error": errMessage,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *scheduledJobRepository) Cancel(
	ctx context.Context,
	id int64,
) error {
	result := r.db.WithContext(ctx).
		Model(&ScheduledJob{}).
		Where("id = ?", id).
		Where(
			"status IN ?",
			[]JobStatus{
				JobStatusPending,
				JobStatusProcessing,
			},
		).
		Update("status", JobStatusCancelled)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *scheduledJobRepository) ResetStaleJobs(
	ctx context.Context,
	before time.Time,
) error {
	return r.db.WithContext(ctx).
		Model(&ScheduledJob{}).
		Where("status = ?", JobStatusProcessing).
		Where("started_at < ?", before).
		Updates(map[string]interface{}{
			"status":     JobStatusPending,
			"started_at": nil,
		}).
		Error
}
