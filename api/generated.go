// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

import (
	"github.com/labstack/echo/v4"
)

// CreateSessionRequest defines model for CreateSessionRequest.
type CreateSessionRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// CreateSessionResponse defines model for CreateSessionResponse.
type CreateSessionResponse struct {
	Token string `json:"token"`
}

// A customer object
type Customer struct {

	// The customer DID
	Id string `json:"id"`

	// Internal name for this customer
	Name string `json:"name"`
}

// CustomersResponse defines model for CustomersResponse.
type CustomersResponse []Customer

// CreateSessionJSONBody defines parameters for CreateSession.
type CreateSessionJSONBody CreateSessionRequest

// CreateSessionJSONRequestBody defines body for CreateSession for application/json ContentType.
type CreateSessionJSONRequestBody CreateSessionJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /web/auth)
	CreateSession(ctx echo.Context) error

	// (GET /web/customers)
	GetCustomers(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// CreateSession converts echo context to params.
func (w *ServerInterfaceWrapper) CreateSession(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateSession(ctx)
	return err
}

// GetCustomers converts echo context to params.
func (w *ServerInterfaceWrapper) GetCustomers(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetCustomers(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/web/auth", wrapper.CreateSession)
	router.GET(baseURL+"/web/customers", wrapper.GetCustomers)

}

