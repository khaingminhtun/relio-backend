package httperror

import (
	"errors"
	"net/http"

	"github.com/khaingminhtun/production-go-api/internal/shared/errorhandler/apperror"
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

// Clean static lookup map to replace the bulky switch block
var codeToStatus = map[apperror.Code]int{
	apperror.CodeUserNotFound:                   http.StatusNotFound,
	apperror.CodeUserAlreadyExists:              http.StatusConflict,
	apperror.CodeInvalidCredentials:             http.StatusUnauthorized,
	apperror.CodeEmailNotVerified:               http.StatusForbidden,
	apperror.CodeInvalidVerifyCode:              http.StatusBadRequest,
	apperror.CodeVerifyCodeExpired:              http.StatusBadRequest,
	apperror.CodeInvalidRequest:                 http.StatusBadRequest,
	apperror.CodeAuthSessionNotFound:            http.StatusNotFound,
	apperror.CodeAccountInactive:                http.StatusForbidden,
	apperror.CodeAuthSessionExpired:             http.StatusUnauthorized,
	apperror.CodeUserInactive:                   http.StatusForbidden,
	apperror.CodeAuthSessionRevoked:             http.StatusUnauthorized,
	apperror.CodeRelationshipNotFound:           http.StatusNotFound,
	apperror.CodeRelationshipMemberNotFound:     http.StatusNotFound,
	apperror.CodeInvitationNotFound:             http.StatusNotFound,
	apperror.CodeInvalidRelationshipName:        http.StatusNotFound,
	apperror.CodeInvalidTimezone:                http.StatusNotFound,
	apperror.CodeInvalidRelationshipType:        http.StatusNotFound,
	apperror.CodeCustomRelationshipTypeRequired: http.StatusNotFound,
	apperror.CodeUnauthorized:                   http.StatusUnauthorized,
	apperror.CodeInvalidAccessToken:             http.StatusUnauthorized,
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
