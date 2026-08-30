package note

import (
	"context"
	"fmt"
	"strings"

	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
)

// ============================================================
// Interface
// ============================================================

type Service interface {
	Create(
		ctx context.Context,
		userID int64,
		req CreateNoteRequest,
	) (*NoteResponse, error)

	List(
		ctx context.Context,
		userID int64,
		filter ListNotesFilter,
	) ([]NoteResponse, error)

	GetByID(
		ctx context.Context,
		userID int64,
		noteID int64,
	) (*NoteResponse, error)

	Update(
		ctx context.Context,
		userID int64,
		noteID int64,
		req UpdateNoteRequest,
	) (*NoteResponse, error)

	Delete(
		ctx context.Context,
		userID int64,
		noteID int64,
	) error
}

// ============================================================
// Implementation
// ============================================================

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// ============================================================
// Create
// ============================================================

func (s *service) Create(
	ctx context.Context,
	userID int64,
	req CreateNoteRequest,
) (*NoteResponse, error) {

	// --------------------------------------------------------
	// Validate title
	// --------------------------------------------------------

	title := strings.TrimSpace(req.Title)

	if title == "" {
		return nil, apperror.New(
			apperror.CodeInvalidRequest,
			"title is required",
			nil,
		)
	}

	// --------------------------------------------------------
	// Validate content
	// --------------------------------------------------------

	content := strings.TrimSpace(req.Content)

	if content == "" {
		return nil, apperror.New(
			apperror.CodeInvalidRequest,
			"content is required",
			nil,
		)
	}

	// --------------------------------------------------------
	// Validate mood (optional)
	// --------------------------------------------------------

	if req.Mood != nil {
		if err := validateMood(*req.Mood); err != nil {
			return nil, err
		}
	}

	// --------------------------------------------------------
	// Create note
	// --------------------------------------------------------

	note := &Note{
		UserID:  userID,
		Title:   title,
		Content: content,
		Mood:    req.Mood,
	}

	if err := s.repo.Create(ctx, note); err != nil {
		return nil, fmt.Errorf("create note: %w", err)
	}

	return toNoteResponse(note), nil
}

// ============================================================
// List
// ============================================================

func (s *service) List(
	ctx context.Context,
	userID int64,
	filter ListNotesFilter,
) ([]NoteResponse, error) {

	// --------------------------------------------------------
	// Validate mood filter when provided
	// --------------------------------------------------------

	if filter.Mood != nil {
		if err := validateMood(*filter.Mood); err != nil {
			return nil, err
		}
	}

	notes, err := s.repo.List(ctx, userID, filter)
	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}

	responses := make([]NoteResponse, 0, len(notes))

	for _, n := range notes {
		responses = append(responses, *toNoteResponse(&n))
	}

	return responses, nil
}

// ============================================================
// GetByID
// ============================================================

func (s *service) GetByID(
	ctx context.Context,
	userID int64,
	noteID int64,
) (*NoteResponse, error) {

	// user_id scoping is enforced inside the repository
	note, err := s.repo.GetByID(ctx, userID, noteID)
	if err != nil {
		return nil, err
	}

	return toNoteResponse(note), nil
}

// ============================================================
// Update
// ============================================================

func (s *service) Update(
	ctx context.Context,
	userID int64,
	noteID int64,
	req UpdateNoteRequest,
) (*NoteResponse, error) {

	// --------------------------------------------------------
	// Fetch note — user_id scoping enforced in repository
	// --------------------------------------------------------

	note, err := s.repo.GetByID(ctx, userID, noteID)
	if err != nil {
		return nil, err
	}

	// --------------------------------------------------------
	// Apply partial updates
	// --------------------------------------------------------

	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)

		if title == "" {
			return nil, apperror.New(
				apperror.CodeInvalidRequest,
				"title cannot be empty",
				nil,
			)
		}

		note.Title = title
	}

	if req.Content != nil {
		content := strings.TrimSpace(*req.Content)

		if content == "" {
			return nil, apperror.New(
				apperror.CodeInvalidRequest,
				"content cannot be empty",
				nil,
			)
		}

		note.Content = content
	}

	if req.Mood != nil {
		if err := validateMood(*req.Mood); err != nil {
			return nil, err
		}

		note.Mood = req.Mood
	}

	if req.IsPinned != nil {
		note.IsPinned = *req.IsPinned
	}

	if req.IsArchived != nil {
		note.IsArchived = *req.IsArchived
	}

	if err := s.repo.Update(ctx, note); err != nil {
		return nil, fmt.Errorf("update note: %w", err)
	}

	return toNoteResponse(note), nil
}

// ============================================================
// Delete
// ============================================================

func (s *service) Delete(
	ctx context.Context,
	userID int64,
	noteID int64,
) error {

	// user_id scoping enforced in repository
	if err := s.repo.Delete(ctx, userID, noteID); err != nil {
		return fmt.Errorf("delete note: %w", err)
	}

	return nil
}

// ============================================================
// Helpers
// ============================================================

func validateMood(mood string) error {
	if _, ok := ValidMoods[mood]; !ok {
		return apperror.New(
			apperror.CodeInvalidRequest,
			"invalid mood value",
			nil,
		)
	}

	return nil
}

func toNoteResponse(n *Note) *NoteResponse {
	return &NoteResponse{
		ID:         n.ID,
		UserID:     n.UserID,
		Title:      n.Title,
		Content:    n.Content,
		Mood:       n.Mood,
		IsPinned:   n.IsPinned,
		IsArchived: n.IsArchived,
		CreatedAt:  n.CreatedAt,
		UpdatedAt:  n.UpdatedAt,
	}
}
