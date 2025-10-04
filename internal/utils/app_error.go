package utils

import "fmt"

type AppError struct {
	ErrorCode     ErrorCode
	Err           error
	CustomMessage *string
}

func (e *AppError) Error() string {
	message := e.ErrorCode.Message
	if e.CustomMessage != nil && *e.CustomMessage != "" {
		message = *e.CustomMessage
	}

	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.ErrorCode.Key, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", e.ErrorCode.Key, message)
}

func (e *AppError) GetDisplayMessage() string {
	if e.CustomMessage != nil && *e.CustomMessage != "" {
		return *e.CustomMessage
	}
	return e.ErrorCode.Message
}

func NewAppError(resp ErrorCode, err error) *AppError {
	errorInfo := GetError(resp)
	return &AppError{
		ErrorCode:     errorInfo,
		Err:           err,
		CustomMessage: nil,
	}
}

func NewAppErrorWithMessage(resp ErrorCode, err error, customMessage string) *AppError {
	errorInfo := GetError(resp)
	return &AppError{
		ErrorCode:     errorInfo,
		Err:           err,
		CustomMessage: &customMessage,
	}
}

func NewAppErrorMessage(resp ErrorCode, customMessage string) *AppError {
	errorInfo := GetError(resp)
	return &AppError{
		ErrorCode:     errorInfo,
		Err:           nil,
		CustomMessage: &customMessage,
	}
}

func (e *AppError) WithMessage(message string) *AppError {
	e.CustomMessage = &message
	return e
}
