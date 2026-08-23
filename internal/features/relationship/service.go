package relationship

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	transaction "github.com/khaingminhtun/production-go-api/internal/shared/dbutils"
	"github.com/khaingminhtun/production-go-api/internal/shared/errorhandler/apperror"
	"gorm.io/gorm"
)

type Service interface {
	Create(
		ctx context.Context,
		userID int64,
		req CreateRelationshipRequest,
	) (*RelationshipResponse, error)

	List(ctx context.Context, userID int64) ([]RelationshipResponse, error)

	GetByID(
		ctx context.Context,
		id int64,
	) (*RelationshipResponse, error)

	Update(
		ctx context.Context,
		id int64,
		req UpdateRelationshipRequest) (*RelationshipResponse, error)

	Delete(
		ctx context.Context,
		id int64) error
}

type service struct {
	relationshipRepo RelationshipRepository
	memberRepo       RelationshipMemberRepository
	txManager        transaction.Manager
}

func NewService(relationshipRepo RelationshipRepository, memberRepo RelationshipMemberRepository, txManager transaction.Manager) Service {
	return &service{
		relationshipRepo: relationshipRepo,
		memberRepo:       memberRepo,
		txManager:        txManager,
	}
}

func (s service) Create(ctx context.Context, userID int64, req CreateRelationshipRequest) (*RelationshipResponse, error) {

	name := strings.TrimSpace(req.Name)
	timezone := strings.TrimSpace(req.Timezone)

	var customType *string

	if req.CustomType != nil {
		value := strings.TrimSpace(*req.CustomType)

		if value != "" {
			customType = &value
		}
	}

	if name == "" {
		return nil, apperror.New(
			apperror.CodeInvalidRelationshipName,
			"relationship name is required",
			nil,
		)
	}

	if timezone == "" {
		return nil, apperror.New(
			apperror.CodeInvalidTimezone,
			"timezone is required",
			nil,
		)
	}

	if !req.Type.IsValid() {
		return nil, apperror.New(
			apperror.CodeInvalidRelationshipType,
			"invalid relationship type",
			nil,
		)
	}

	if req.Type == RelationshipTypeOther {

		if customType == nil {
			return nil, apperror.New(
				apperror.CodeCustomRelationshipTypeRequired,
				"custom type is required when relationship type is other",
				nil,
			)
		}

	} else {
		// custom_type is only allowed when type = other.
		customType = nil
	}

	relationship := &Relationship{
		Name:        name,
		Type:        req.Type,
		CustomType:  customType,
		Description: req.Description,
		StartDate:   req.StartDate,
		Timezone:    timezone,
		CreatedBy:   userID,
	}

	now := time.Now()
	member := &RelationshipMember{
		RelationshipID: relationship.ID,
		UserID:         userID,
		Role:           RelationshipMemberRoleOwner,
		Status:         RelationshipMemberStatusActive,
		JoinedAt:       &now,
	}

	// ============================================================
	// Create relationship + owner member in one transaction
	// ============================================================

	err := s.txManager.WithinTransaction(
		ctx,
		func(ctx context.Context) error {

			// --------------------------------------------------------
			// Create relationship
			// --------------------------------------------------------

			if err := s.relationshipRepo.Create(
				ctx,
				relationship,
			); err != nil {
				return fmt.Errorf(
					"create relationship: %w",
					err,
				)
			}

			// Relationship ID is generated after Create().
			member.RelationshipID = relationship.ID

			// --------------------------------------------------------
			// Create owner member
			// --------------------------------------------------------

			if err := s.memberRepo.Create(
				ctx,
				member,
			); err != nil {
				return fmt.Errorf(
					"create relationship member: %w",
					err,
				)
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	// ============================================================
	// Response
	// ============================================================

	return toRelationshipResponse(relationship), nil

}

func (s service) List(ctx context.Context, userID int64) ([]RelationshipResponse, error) {
	relationships, err := s.relationshipRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, apperror.New(
			apperror.CodeRelationshipNotFound,
			"relationships not found ",
			nil,
		)
	}

	responses := make([]RelationshipResponse, 0, len(relationships))

	for _, relationship := range relationships {
		responses = append(responses, RelationshipResponse{
			ID:          relationship.ID,
			Name:        relationship.Name,
			Type:        relationship.Type,
			CustomType:  relationship.CustomType,
			Description: relationship.Description,
			StartDate:   relationship.StartDate,
			Timezone:    relationship.Timezone,
			CreatedBy:   relationship.CreatedBy,
			CreatedAt:   relationship.CreatedAt,
			UpdatedAt:   relationship.UpdatedAt,
		})
	}

	return responses, nil
}

func (s *service) GetByID(
	ctx context.Context,
	id int64,
) (*RelationshipResponse, error) {
	relationship, err := s.relationshipRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeRelationshipNotFound,
				"relationship not found",
				err,
			)
		}

		return nil, fmt.Errorf("get relationship by id: %w", err)
	}

	return &RelationshipResponse{
		ID:          relationship.ID,
		Name:        relationship.Name,
		Type:        relationship.Type,
		CustomType:  relationship.CustomType,
		Description: relationship.Description,
		StartDate:   relationship.StartDate,
		Timezone:    relationship.Timezone,
		CreatedBy:   relationship.CreatedBy,
		CreatedAt:   relationship.CreatedAt,
		UpdatedAt:   relationship.UpdatedAt,
	}, nil
}

func (s *service) Update(
	ctx context.Context,
	id int64,
	req UpdateRelationshipRequest,
) (*RelationshipResponse, error) {
	relationship, err := s.relationshipRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeRelationshipNotFound,
				"relationship not found",
				err,
			)
		}

		return nil, fmt.Errorf("get relationship for update: %w", err)
	}

	if req.Name != nil {
		relationship.Name = *req.Name
	}

	if req.Type != nil {
		relationship.Type = *req.Type
	}

	if req.CustomType != nil {
		relationship.CustomType = req.CustomType
	}

	if req.Description != nil {
		relationship.Description = req.Description
	}

	if req.StartDate != nil {
		relationship.StartDate = req.StartDate
	}

	if req.Timezone != nil {
		relationship.Timezone = *req.Timezone
	}

	if err := s.relationshipRepo.Update(ctx, relationship); err != nil {
		return nil, fmt.Errorf("update relationship: %w", err)
	}

	return &RelationshipResponse{
		ID:          relationship.ID,
		Name:        relationship.Name,
		Type:        relationship.Type,
		CustomType:  relationship.CustomType,
		Description: relationship.Description,
		StartDate:   relationship.StartDate,
		Timezone:    relationship.Timezone,
		CreatedBy:   relationship.CreatedBy,
		CreatedAt:   relationship.CreatedAt,
		UpdatedAt:   relationship.UpdatedAt,
	}, nil
}

func (s *service) Delete(
	ctx context.Context,
	id int64,
) error {
	err := s.relationshipRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.New(
				apperror.CodeRelationshipNotFound,
				"relationship not found",
				err,
			)
		}

		return fmt.Errorf("delete relationship: %w", err)
	}

	return nil
}

func toRelationshipResponse(
	relationship *Relationship,
) *RelationshipResponse {
	return &RelationshipResponse{
		ID:          relationship.ID,
		Name:        relationship.Name,
		Type:        relationship.Type,
		CustomType:  relationship.CustomType,
		Description: relationship.Description,
		StartDate:   relationship.StartDate,
		Timezone:    relationship.Timezone,
		CreatedBy:   relationship.CreatedBy,
		CreatedAt:   relationship.CreatedAt,
		UpdatedAt:   relationship.UpdatedAt,
	}
}
