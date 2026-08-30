package memory

import "time"

// ============================================================
// Request DTOs
// ============================================================

type CreateMemoryRequest struct {
	Title      string     `json:"title"`
	Content    *string    `json:"content"`
	MemoryDate *time.Time `json:"memory_date"`
}

type UpdateMemoryRequest struct {
	Title      *string    `json:"title"`
	Content    *string    `json:"content"`
	MemoryDate *time.Time `json:"memory_date"`
}

// ============================================================
// Response DTO
// ============================================================

type MemoryResponse struct {
	ID             int64      `json:"id"`
	RelationshipID int64      `json:"relationship_id"`
	CreatedBy      int64      `json:"created_by"`
	Title          string     `json:"title"`
	Content        *string    `json:"content,omitempty"`
	MemoryDate     *time.Time `json:"memory_date,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
