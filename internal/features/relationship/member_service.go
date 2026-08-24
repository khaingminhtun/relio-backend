package relationship

import (
	"context"

	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
)

type RelationshipMemberService interface {
	ListByRelationshipID(
		ctx context.Context,
		relationshipID int64,
	) ([]RelationshipMemberResponse, error)

	GetByID(
		ctx context.Context,
		id int64,
	) (*RelationshipMemberResponse, error)

	GetByRelationshipAndUser(
		ctx context.Context,
		relationshipID int64,
		userID int64,
	) (*RelationshipMemberResponse, error)

	UpdateMember(
		ctx context.Context,
		memberID int64,
		req UpdateMemberRequest,
	) (*RelationshipMemberResponse, error)

	RemoveMember(
		ctx context.Context,
		relationshipID int64,
		userID int64,
	) error
}

type relationshipMemberService struct {
	memberRepo RelationshipMemberRepository
}

func NewRelationshipMemberService(
	memberRepo RelationshipMemberRepository,
) RelationshipMemberService {
	return &relationshipMemberService{
		memberRepo: memberRepo,
	}
}

func (s relationshipMemberService) ListByRelationshipID(
	ctx context.Context,
	relationshipID int64,
) ([]RelationshipMemberResponse, error) {

	members, err := s.memberRepo.ListByRelationshipID(
		ctx,
		relationshipID,
	)
	if err != nil {
		return nil, apperror.New(
			apperror.CodeRelationshipMemberNotFound,
			"relationship members not found",
			nil,
		)
	}

	responses := make(
		[]RelationshipMemberResponse,
		0,
		len(members),
	)

	for _, member := range members {
		responses = append(responses, *toMemberResponse(&member))
	}

	return responses, nil
}

func (s relationshipMemberService) GetByID(
	ctx context.Context,
	id int64,
) (*RelationshipMemberResponse, error) {
	member, err := s.memberRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.New(
			apperror.CodeRelationshipMemberNotFound,
			"relationship members not found",
			nil,
		)
	}

	return toMemberResponse(member), nil
}

func (s relationshipMemberService) GetByRelationshipAndUser(
	ctx context.Context,
	relationshipID int64,
	userID int64,
) (*RelationshipMemberResponse, error) {
	member, err := s.memberRepo.GetByRelationshipAndUser(
		ctx,
		relationshipID,
		userID,
	)
	if err != nil {
		return nil, apperror.New(
			apperror.CodeRelationshipMemberNotFound,
			"relationship members not found",
			nil,
		)
	}

	return toMemberResponse(member), nil
}

func (s relationshipMemberService) UpdateMember(
	ctx context.Context,
	memberID int64,
	req UpdateMemberRequest,
) (*RelationshipMemberResponse, error) {
	member, err := s.memberRepo.GetByID(ctx, memberID)
	if err != nil {
		return nil, apperror.New(
			apperror.CodeRelationshipMemberNotFound,
			"relationship members not found",
			nil,
		)
	}

	if req.Role != nil {
		member.Role = *req.Role
	}

	if req.Status != nil {
		member.Status = *req.Status
	}

	if err := s.memberRepo.Update(ctx, member); err != nil {
		return nil, apperror.New(
			apperror.CodeInvalidRequest,
			"failed to update relationship member",
			err,
		)
	}

	return toMemberResponse(member), nil
}

func (s relationshipMemberService) RemoveMember(
	ctx context.Context,
	relationshipID int64,
	userID int64,
) error {
	return s.memberRepo.Delete(ctx, relationshipID, userID)
}

func toMemberResponse(member *RelationshipMember) *RelationshipMemberResponse {
	return &RelationshipMemberResponse{
		ID:             member.ID,
		RelationshipID: member.RelationshipID,
		UserID:         member.UserID,
		Role:           member.Role,
		Status:         member.Status,
		JoinedAt:       member.JoinedAt,
		CreatedAt:      member.CreatedAt,
		UpdatedAt:      member.UpdatedAt,
	}
}
