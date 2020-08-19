package apperror

import (
	"fmt"
	"net/http"
)

// AuthError defines auth failed error.
type AuthError struct {
}

// Error implements error interface.
func (e *AuthError) Error() string {
	return "please login"
}

// Code implements HTTPError interface.
func (e *AuthError) Code() int {
	return http.StatusUnauthorized
}

// BaseInfo describes that an user attmpts to access an endpoint.
type BaseInfo struct {
	UserEmail string
	Method    string
	URI       string
}

// SetBaseInfo reads info from appctx.
func (info *BaseInfo) SetBaseInfo(email, method, URI string) {
	info.UserEmail = email
	info.Method = method
	info.URI = URI
}

// GeneralError defines general HTTP error.
type GeneralError struct {
	BaseInfo
	Err error
}

// Error implements error interface.
func (e *GeneralError) Error() string {
	return fmt.Sprintf("user(%v) sent req to method(%v) uri(%v) failed(%v)",
		e.UserEmail, e.Method, e.URI, e.Err)
}

// Code implements HTTPError interface.
func (e *GeneralError) Code() int {
	return http.StatusInternalServerError
}

// BadParamError defines error that parameter is illegal.
type BadParamError struct {
	BaseInfo
	ParsingError error
}

// Error implements error interface.
func (e *BadParamError) Error() string {
	return fmt.Sprintf("user(%v) sent req to method(%v) uri(%v) failed(%v)",
		e.UserEmail, e.Method, e.URI, e.ParsingError)
}

// Code implements HTTPError interface.
func (e *BadParamError) Code() int {
	return http.StatusBadRequest
}

// BackendFailedError defines error that backend service failed.
type BackendFailedError struct {
	BaseInfo
	Action        string
	ServiceType   string
	ServiceAddr   string
	InternalError error
}

// Error implements error interface.
func (e *BackendFailedError) Error() string {
	return fmt.Sprintf("user(%v) from method(%v) uri(%v) accessed backend "+
		"type(%v) addr(%v) do action(%v) failed(%v)", e.UserEmail, e.Method, e.URI,
		e.ServiceType, e.ServiceAddr, e.Action, e.InternalError)
}

// Code implements HTTPError interface.
func (e *BackendFailedError) Code() int {
	return http.StatusGatewayTimeout
}

// DatabaseFailedError defines error that database returns error.
type DatabaseFailedError struct {
	BaseInfo
	DBError error
}

// Error implements error interface.
func (e *DatabaseFailedError) Error() string {
	return fmt.Sprintf("user(%v) from method(%v) uri(%v) accessed db(%v)",
		e.UserEmail, e.Method, e.URI, e.DBError)
}

// Code implements HTTPError interface.
func (e *DatabaseFailedError) Code() int {
	return http.StatusInternalServerError
}
