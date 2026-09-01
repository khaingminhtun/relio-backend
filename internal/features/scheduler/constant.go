package scheduler

type JobStatus string

type JobType string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
	JobStatusCancelled  JobStatus = "cancelled"
)

const (
	JobTypeEventReminder JobType = "event_reminder"
)

func (s JobStatus) IsValid() bool {
	switch s {
	case JobStatusPending,
		JobStatusProcessing,
		JobStatusCompleted,
		JobStatusFailed,
		JobStatusCancelled:
		return true
	default:
		return false
	}
}

func (t JobType) IsValid() bool {
	switch t {
	case JobTypeEventReminder:
		return true
	default:
		return false
	}
}
