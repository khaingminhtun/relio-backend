package relationship

import "context"

type RelationshipRepository interface {
	Create(
		ctx context.Context,
		relationship *Relationship,
	) error

	GetByID(
		ctx context.Context,
		id int64,
	) (*Relationship, error)

	ListByUserID(
		ctx context.Context,
		userID int64,
	) ([]Relationship, error)

	Update(
		ctx context.Context,
		relationship *Relationship,
	) error

	Delete(
		ctx context.Context,
		id int64,
	) error
}

type RelationshipMemberRepository interface {
	Create(
		ctx context.Context,
		member *RelationshipMember,
	) error

	GetByID(
		ctx context.Context,
		id int64,
	) (*RelationshipMember, error)

	GetByRelationshipAndUser(
		ctx context.Context,
		relationshipID int64,
		userID int64,
	) (*RelationshipMember, error)

	ListByRelationshipID(
		ctx context.Context,
		relationshipID int64,
	) ([]RelationshipMember, error)

	Update(
		ctx context.Context,
		member *RelationshipMember,
	) error

	Delete(
		ctx context.Context,
		relationshipID int64,
		userID int64,
	) error

	Exists(
		ctx context.Context,
		relationshipID int64,
		userID int64,
	) (bool, error)
}

type InvitationRepository interface {
	Create(
		ctx context.Context,
		invitation *Invitation,
	) error

	GetByID(
		ctx context.Context,
		id int64,
	) (*Invitation, error)

	GetByTokenHash(
		ctx context.Context,
		tokenHash string,
	) (*Invitation, error)

	ListByRelationshipID(
		ctx context.Context,
		relationshipID int64,
	) ([]Invitation, error)

	FindPendingByEmail(
		ctx context.Context,
		relationshipID int64,
		email string,
	) (*Invitation, error)

	Update(
		ctx context.Context,
		invitation *Invitation,
	) error

	Delete(
		ctx context.Context,
		id int64,
	) error
}
