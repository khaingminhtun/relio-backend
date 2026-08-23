package relationship

type RelationshipType string

const (
	RelationshipTypeFamily     RelationshipType = "family"
	RelationshipTypeFriendship RelationshipType = "friendship"
	RelationshipTypeCouple     RelationshipType = "couple"
	RelationshipTypeOther      RelationshipType = "other"
)

type RelationshipMemberRole string

const (
	RelationshipMemberRoleOwner  RelationshipMemberRole = "owner"
	RelationshipMemberRoleAdmin  RelationshipMemberRole = "admin"
	RelationshipMemberRoleMember RelationshipMemberRole = "member"
)

type RelationshipMemberStatus string

const (
	RelationshipMemberStatusActive   RelationshipMemberStatus = "active"
	RelationshipMemberStatusInactive RelationshipMemberStatus = "inactive"
	RelationshipMemberStatusLeft     RelationshipMemberStatus = "left"
	RelationshipMemberStatusRemoved  RelationshipMemberStatus = "removed"
)

type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusDeclined InvitationStatus = "declined"
	InvitationStatusExpired  InvitationStatus = "expired"
	InvitationStatusCanceled InvitationStatus = "canceled"
)
