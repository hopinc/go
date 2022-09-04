package structures

import "strconv"

// NotFound is sent when the object is not found. THe string is the error message from the API.
type NotFound string

// Error implements the error interface.
func (n NotFound) Error() string { return (string)(n) }

// ServerError is sent when Hop encounters a internal server error that happens to fit the error schema.
type ServerError string

// Error implements the error interface.
func (s ServerError) Error() string { return (string)(s) }

// NotAuthorized is sent when the user is not authorized. THe string is the error message from the API.
type NotAuthorized string

// Error implements the error interface.
func (n NotAuthorized) Error() string { return (string)(n) }

// BadRequest is sent in the event of a 400. It is a special case since it encompasses all user request errors.
type BadRequest struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// Error implements the error interface.
func (b BadRequest) Error() string {
	return b.Code + ": " + b.Message
}

// UnknownServerError refers to a server error where the cause is unknown.
type UnknownServerError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Code       string `json:"code"`
}

// Error implements the error interface.
func (u UnknownServerError) Error() string {
	return "status code " + strconv.Itoa(u.StatusCode) + " (" + u.Code + "): " + u.Message
}
