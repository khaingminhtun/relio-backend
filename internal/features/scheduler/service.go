package scheduler

import "context"

type SchedulerService interface {
	Run(ctx context.Context) error
	EnqueueDueJobs(ctx context.Context) error
	RecoverStaleJobs(ctx context.Context) error
}
