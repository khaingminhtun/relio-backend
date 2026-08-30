package memory

import (
	"context"
	"fmt"
	"strings"

	"github.com/khaingminhtun/relio-backend/internal/features/relationship"
	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
)

// ============================================================
// Interface
// ============================================================

type MemoryService interface {
	Create(
		ctx context.Context,
		userID int64,
		relationshipID int64,
		req CreateMemoryRequest,
	) (*MemoryResponse, error)

	List(
		ctx context.Context,
		userID int64,
		relationshipID int64,
	) ([]MemoryResponse, error)

	GetByID(
		ctx context.Context,
		userID int64,
		relationshipID int64,
		memoryID int64,
	) (*MemoryResponse, error)

	Update(
		ctx context.Context,
		userID int64,
		relationshipID int64,
		memoryID int64,
		req UpdateMemoryRequest,
	) (*MemoryResponse, error)

	Delete(
		ctx context.Context,
		userID int64,
		relationshipID int64,
		memoryID int64,
	) error
}

// ============================================================
// Implementation
// ============================================================

type memoryService struct {
	memoryRepo relationship.RelationshipMemberRepository
	memberRepo relationship.RelationshipMemberRepository
	repo       MemoryRepository
}

func NewMemoryService(
	repo MemoryRepository,
	memberRepo relationship.RelationshipMemberRepository,
) MemoryService {
	return &memoryService{
		repo:       repo,
		memberRepo: memberRepo,
	}
}

// ============================================================
// Create
// ============================================================

func (s *memoryService) Create(
	ctx context.Context,
	userID int64,
	relationshipID int64,
	req CreateMemoryRequest,
) (*MemoryResponse, error) {

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
	// Check active membership
	// --------------------------------------------------------

	if err := s.requireActiveMember(ctx, relationshipID, userID); err != nil {
		return nil, err
	}

	// --------------------------------------------------------
	// Create memory
	// --------------------------------------------------------

	memory := &Memory{
		RelationshipID: relationshipID,
		CreatedBy:      userID,
		Title:          title,
		Content:        req.Content,
		MemoryDate:     req.MemoryDate,
	}

	if err := s.repo.Create(ctx, memory); err != nil {
		return nil, fmt.Errorf("create memory: %w", err)
	}

	return toMemoryResponse(memory), nil
}

// ============================================================
// List
// ============================================================

func (s *memoryService) List(
	ctx context.Context,
	userID int64,
	relationshipID int64,
) ([]MemoryResponse, error) {

	// --------------------------------------------------------
	// Check active membership
	// --------------------------------------------------------

	if err := s.requireActiveMember(ctx, relationshipID, userID); err != nil {
		return nil, err
	}

	// --------------------------------------------------------
	// Fetch memories
	// --------------------------------------------------------

	memories, err := s.repo.ListByRelationshipID(ctx, relationshipID)
	if err != nil {
		return nil, fmt.Errorf("list memories: %w", err)
	}

	responses := make([]MemoryResponse, 0, len(memories))

	for _, m := range memories {
		responses = append(responses, *toMemoryResponse(&m))
	}

	return responses, nil
}

// ============================================================
// GetByID
// ============================================================

func (s *memoryService) GetByID(
	ctx context.Context,
	userID int64,
	relationshipID int64,
	memoryID int64,
) (*MemoryResponse, error) {

	// --------------------------------------------------------
	// Check active membership
	// --------------------------------------------------------

	if err := s.requireActiveMember(ctx, relationshipID, userID); err != nil {
		return nil, err
	}

	// --------------------------------------------------------
	// Fetch memory (relationship_id enforced inside repo)
	// --------------------------------------------------------

	memory, err := s.repo.GetByID(ctx, relationshipID, memoryID)
	if err != nil {
		return nil, err
	}

	return toMemoryResponse(memory), nil
}

// ============================================================
// Update
// ============================================================

func (s *memoryService) Update(
	ctx context.Context,
	userID int64,
	relationshipID int64,
	memoryID int64,
	req UpdateMemoryRequest,
) (*MemoryResponse, error) {

	// --------------------------------------------------------
	// Fetch memory (relationship_id enforced inside repo)
	// --------------------------------------------------------

	memory, err := s.repo.GetByID(ctx, relationshipID, memoryID)
	if err != nil {
		return nil, err
	}

	// --------------------------------------------------------
	// Authorization: creator OR owner/admin of the relationship
	// --------------------------------------------------------

	if err := s.requireCreatorOrPrivileged(ctx, relationshipID, userID, memory.CreatedBy); err != nil {
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

		memory.Title = title
	}

	if req.Content != nil {
		memory.Content = req.Content
	}

	if req.MemoryDate != nil {
		memory.MemoryDate = req.MemoryDate
	}

	if err := s.repo.Update(ctx, memory); err != nil {
		return nil, fmt.Errorf("update memory: %w", err)
	}

	return toMemoryResponse(memory), nil
}

// ============================================================
// Delete
// ============================================================

func (s *memoryService) Delete(
	ctx context.Context,
	userID int64,
	relationshipID int64,
	memoryID int64,
) error {

	// --------------------------------------------------------
	// Fetch memory (relationship_id enforced inside repo)
	// --------------------------------------------------------

	memory, err := s.repo.GetByID(ctx, relationshipID, memoryID)
	if err != nil {
		return err
	}

	// --------------------------------------------------------
	// Authorization: creator OR owner/admin of the relationship
	// --------------------------------------------------------

	if err := s.requireCreatorOrPrivileged(ctx, relationshipID, userID, memory.CreatedBy); err != nil {
		return err
	}

	// --------------------------------------------------------
	// Soft delete
	// --------------------------------------------------------

	if err := s.repo.Delete(ctx, relationshipID, memoryID); err != nil {
		return fmt.Errorf("delete memory: %w", err)
	}

	return nil
}

// ============================================================
// Authorization helpers
// ============================================================

// requireActiveMember returns an error if the user is not an active
// member of the relationship.
func (s *memoryService) requireActiveMember(
	ctx context.Context,
	relationshipID int64,
	userID int64,
) error {
	exists, err := s.memberRepo.Exists(ctx, relationshipID, userID)
	if err != nil {
		return fmt.Errorf("check membership: %w", err)
	}

	if !exists {
		return apperror.New(
			apperror.CodeInvalidRelationshipMember,
			"user is not an active member of this relationship",
			nil,
		)
	}

	return nil
}

// requireCreatorOrPrivileged returns an error unless the user is the
// memory creator OR is an owner/admin of the relationship.
func (s *memoryService) requireCreatorOrPrivileged(
	ctx context.Context,
	relationshipID int64,
	userID int64,
	createdBy int64,
) error {
	// Creator always has access.
	if userID == createdBy {
		return nil
	}

	// Check if the user is an owner/admin of the relationship.
	member, err := s.memberRepo.GetByRelationshipAndUser(ctx, relationshipID, userID)
	if err != nil {
		// User is not a member at all.
		return apperror.New(
			apperror.CodeMemoryForbidden,
			"you do not have permission to modify this memory",
			nil,
		)
	}

	if member.Role == relationship.RelationshipMemberRoleOwner ||
		member.Role == relationship.RelationshipMemberRoleAdmin {
		return nil
	}

	return apperror.New(
		apperror.CodeMemoryForbidden,
		"you do not have permission to modify this memory",
		nil,
	)
}

// ============================================================
// Mapper
// ============================================================

func toMemoryResponse(m *Memory) *MemoryResponse {
	return &MemoryResponse{
		ID:             m.ID,
		RelationshipID: m.RelationshipID,
		CreatedBy:      m.CreatedBy,
		Title:          m.Title,
		Content:        m.Content,
		MemoryDate:     m.MemoryDate,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
