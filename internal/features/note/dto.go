package note

import "time"

// ============================================================
// List filter
// ============================================================

type ListNotesFilter struct {
	Mood       *string
	IsPinned   *bool
	IsArchived *bool
	Search     *string
}

// ============================================================
// Request DTOs
// ============================================================

type CreateNoteRequest struct {
	Title   string  `json:"title"`
	Content string  `json:"content"`
	Mood    *string `json:"mood"`
}

type UpdateNoteRequest struct {
	Title      *string `json:"title"`
	Content    *string `json:"content"`
	Mood       *string `json:"mood"`
	IsPinned   *bool   `json:"is_pinned"`
	IsArchived *bool   `json:"is_archived"`
}

// ============================================================
// Response DTO
// ============================================================

type NoteResponse struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	Mood       *string    `json:"mood,omitempty"`
	IsPinned   bool       `json:"is_pinned"`
	IsArchived bool       `json:"is_archived"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
