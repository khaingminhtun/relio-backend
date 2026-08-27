package httperror

import (
	"errors"
	"net/http"

	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
)

type Error struct {
	Status  int
	Code    string
	Message string
	Err     error
}

func New(status int, code, message string, err error) *Error {
	return &Error{
		Status:  status,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error() // Improved string tracking
	}

	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func FromError(err error) *Error {
	if err == nil {
		return nil
	}

	// Already an HTTP error.
	var httpErr *Error
	if errors.As(err, &httpErr) {
		return httpErr
	}

	// Application/business error.
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		return fromAppError(appErr)
	}

	// Unknown/System error (Database crashes, network connection loss, etc.)
	return New(
		http.StatusInternalServerError,
		"INTERNAL_SERVER_ERROR",
		"internal server error",
		err,
	)
}

var codeToStatus = map[apperror.Code]int{
	// User
	apperror.CodeUserNotFound:      http.StatusNotFound,
	apperror.CodeUserAlreadyExists: http.StatusConflict,
	apperror.CodeUserInactive:      http.StatusForbidden,
	apperror.CodeAccountInactive:   http.StatusForbidden,

	// Authentication
	apperror.CodeInvalidCredentials:  http.StatusUnauthorized,
	apperror.CodeEmailNotVerified:    http.StatusForbidden,
	apperror.CodeAuthSessionNotFound: http.StatusNotFound,
	apperror.CodeAuthSessionExpired:  http.StatusUnauthorized,
	apperror.CodeAuthSessionRevoked:  http.StatusUnauthorized,
	apperror.CodeUnauthorized:        http.StatusUnauthorized,
	apperror.CodeInvalidAccessToken:  http.StatusUnauthorized,

	// Verification
	apperror.CodeInvalidVerifyCode: http.StatusBadRequest,
	apperror.CodeVerifyCodeExpired: http.StatusBadRequest,

	// Request / Validation
	apperror.CodeInvalidRequest:                 http.StatusBadRequest,
	apperror.CodeInvalidRelationshipName:        http.StatusBadRequest,
	apperror.CodeInvalidTimezone:                http.StatusBadRequest,
	apperror.CodeInvalidRelationshipType:        http.StatusBadRequest,
	apperror.CodeCustomRelationshipTypeRequired: http.StatusBadRequest,
	apperror.CodeInvalidRelationshipMember:      http.StatusBadRequest,

	// Relationship
	apperror.CodeRelationshipNotFound:       http.StatusNotFound,
	apperror.CodeRelationshipMemberNotFound: http.StatusNotFound,

	// Profile
	apperror.CodeUserProfileNotFound: http.StatusNotFound,

	// Invitation
	apperror.CodeInvitationNotFound:      http.StatusNotFound,
	apperror.CodeInvitationExpired:       http.StatusGone,
	apperror.CodeInvitationInvalid:       http.StatusBadRequest,
	apperror.CodeInvitationEmailMismatch: http.StatusForbidden,
}

func fromAppError(err *apperror.Error) *Error {
	status, exists := codeToStatus[err.Code]
	if !exists {
		return New(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			"internal server error",
			err,
		)
	}

	return New(status, string(err.Code), err.Message, err)
}
