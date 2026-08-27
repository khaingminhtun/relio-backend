package relationship

import "time"

type CreateRelationshipRequest struct {
	Name        string           `json:"name"`
	Type        RelationshipType `json:"type"`
	CustomType  *string          `json:"custom_type,omitempty"`
	Description *string          `json:"description,omitempty"`
	StartDate   *time.Time       `json:"start_date,omitempty"`
	Timezone    string           `json:"timezone"`
}

type RelationshipResponse struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	Type        RelationshipType `json:"type"`
	CustomType  *string          `json:"custom_type,omitempty"`
	Description *string          `json:"description,omitempty"`
	StartDate   *time.Time       `json:"start_date,omitempty"`
	Timezone    string           `json:"timezone"`
	CreatedBy   int64            `json:"created_by"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type UpdateRelationshipRequest struct {
	Name        *string           `json:"name"`
	Type        *RelationshipType `json:"type"`
	CustomType  *string           `json:"custom_type"`
	Description *string           `json:"description"`
	StartDate   *time.Time        `json:"start_date"`
	Timezone    *string           `json:"timezone"`
}

// relationship member
type RelationshipMemberResponse struct {
	ID             int64                    `json:"id"`
	RelationshipID int64                    `json:"relationship_id"`
	UserID         int64                    `json:"user_id"`
	Role           RelationshipMemberRole   `json:"role"`
	Status         RelationshipMemberStatus `json:"status"`
	JoinedAt       *time.Time               `json:"joined_at,omitempty"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
}

type UpdateMemberRequest struct {
	Role   *RelationshipMemberRole   `json:"role"`
	Status *RelationshipMemberStatus `json:"status"`
}

type ExternalInvitationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ExternalInvitationResponse struct {
	ID             int64      `json:"id"`
	RelationshipID int64      `json:"relationship_id"`
	Status         string     `json:"status"`
	InviteURL      string     `json:"invite_url"`
	ExpiresAt      time.Time  `json:"expires_at"`
	AcceptedAt     *time.Time `json:"accepted_at"`
}
