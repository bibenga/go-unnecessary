// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.2 DO NOT EDIT.
package server

const (
	ApiKeyAuthScopes = "ApiKeyAuth.Scopes"
	BasicAuthScopes  = "BasicAuth.Scopes"
)

// Application defines model for Application.
type Application struct {
	DictPlatform   *DictPlatform `json:"dictPlatform,omitempty"`
	DictPlatformId int32         `json:"dictPlatformId"`
	DisplayName    *string       `json:"displayName,omitempty"`
	Id             int64         `json:"id"`
	Name           string        `json:"name"`
}

// DictPlatform defines model for DictPlatform.
type DictPlatform struct {
	DisplayName *string `json:"displayName,omitempty"`
	Id          int32   `json:"id"`
	Name        string  `json:"name"`
}

// Error defines model for Error.
type Error struct {
	// Code Error code
	Code *int32 `json:"code,omitempty"`

	// Error Error description
	Error *string `json:"error,omitempty"`

	// Message Error message
	Message string `json:"message"`
}

// ExtraError defines model for ExtraError.
type ExtraError struct {
	// Code Error code
	Code *int32 `json:"code,omitempty"`

	// Error Error description
	Error *string `json:"error,omitempty"`

	// Message Error message
	Message string `json:"message"`

	// RootCause Error location
	RootCause *string `json:"rootCause,omitempty"`
}

// GetStatusV1 defines model for GetStatusV1.
type GetStatusV1 struct {
	Cpu    *int64 `json:"cpu,omitempty"`
	Status string `json:"status"`
}

// AuthenticationError defines model for AuthenticationError.
type AuthenticationError = Error

// BadRequest defines model for BadRequest.
type BadRequest = Error

// PermissionDenid defines model for PermissionDenid.
type PermissionDenid = Error

// GetStatusV1Params defines parameters for GetStatusV1.
type GetStatusV1Params struct {
	Q         *string `form:"q,omitempty" json:"q,omitempty"`
	IsFull    *bool   `form:"IsFull,omitempty" json:"IsFull,omitempty"`
	XPage     *int32  `json:"X-Page,omitempty"`
	XPageSize *int32  `json:"X-Page-Size,omitempty"`
}

// SetStatusV1JSONRequestBody defines body for SetStatusV1 for application/json ContentType.
type SetStatusV1JSONRequestBody = GetStatusV1