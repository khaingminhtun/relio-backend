package apperror

import "errors"

type Code string

const (
	CodeUserNotFound                   Code = "USER_NOT_FOUND"
	CodeUserAlreadyExists              Code = "USER_ALREADY_EXISTS"
	CodeInvalidCredentials             Code = "INVALID_CREDENTIALS"
	CodeEmailNotVerified               Code = "EMAIL_NOT_VERIFIED"
	CodeInvalidVerifyCode              Code = "INVALID_VERIFICATION_CODE"
	CodeVerifyCodeExpired              Code = "VERIFICATION_CODE_EXPIRED"
	CodeInvalidRequest                 Code = "INVALID_REQUEST"
	CodeAuthSessionNotFound            Code = "AUTH_SESSION_NOT_FOUND"
	CodeAccountInactive                Code = "ACCOUNT_INACTIVE"
	CodeInvalidRefreshToken            Code = "REFRESH_TOKEN_INVALID"
	CodeAuthSessionExpired             Code = "AUTH_SESSION_EXPIRED"
	CodeUserInactive                   Code = "USER_INACTIVE"
	CodeAuthSessionRevoked             Code = "AUTH_SESSION_REVOKED"
	CodeRelationshipNotFound           Code = "RELATIONSHIP_NOT_FOUND"
	CodeRelationshipMemberNotFound     Code = "RELATIONSHIP_MEMBER_NOT_FOUND"
	CodeInvitationNotFound             Code = "INVITATION_NOT_FOUND"
	CodeInvalidRelationshipName        Code = "INVALID_RELATIONSHIP_NAME"
	CodeInvalidTimezone                Code = "INVALID_TIME_ZONE"
	CodeInvalidRelationshipType        Code = "INVALID_RELATIONSHIP_TYPE"
	CodeCustomRelationshipTypeRequired Code = "CUSTOM_RELATIONSHIP_TYPE_REQUIRED"
	CodeUnauthorized                   Code = "UNAUTHORIZED"
	CodeInvalidAccessToken             Code = "INVALID_ACCESS_TOKEN"
)

type Error struct {
	Code    Code
	Message string
	Err     error
}

func New(code Code, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func Is(err error, code Code) bool {
	var appErr *Error

	if !errors.As(err, &appErr) {
		return false
	}

	return appErr.Code == code
}
