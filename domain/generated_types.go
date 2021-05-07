// Package domain provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package domain

// ConnectCustomerRequest defines model for ConnectCustomerRequest.
type ConnectCustomerRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	// Used for issueing the NutsOrgCredential
	Town *string `json:"town,omitempty"`
}

// CreateSessionRequest defines model for CreateSessionRequest.
type CreateSessionRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// CreateSessionResponse defines model for CreateSessionResponse.
type CreateSessionResponse struct {
	Token string `json:"token"`
}

// Customer defines model for Customer.
type Customer struct {

	// If a VC has been issued for this customer.
	Active bool `json:"active"`

	// The customer DID.
	Did string `json:"did"`

	// The email domain of the care providers employees, required for logging in.
	Domain *string `json:"domain,omitempty"`

	// The internal customer ID.
	Id string `json:"id"`

	// Internal name for this customer.
	Name string `json:"name"`

	// Locality for this customer.
	Town *string `json:"town,omitempty"`
}

// CustomersResponse defines model for CustomersResponse.
type CustomersResponse []Customer

// ServiceProvider defines model for ServiceProvider.
type ServiceProvider struct {

	// Email address available for other service providers in the network for getting support
	Email string `json:"email"`

	// The DID of the service provider
	Id string `json:"id"`

	// The name of the service provider
	Name string `json:"name"`

	// Number available for other service providers in the network to call in case of emergency
	Phone string `json:"phone"`

	// Publicly reachable website address of the service provider
	Website string `json:"website"`
}

// CreateSessionJSONBody defines parameters for CreateSession.
type CreateSessionJSONBody CreateSessionRequest

// ConnectCustomerJSONBody defines parameters for ConnectCustomer.
type ConnectCustomerJSONBody ConnectCustomerRequest

// UpdateCustomerJSONBody defines parameters for UpdateCustomer.
type UpdateCustomerJSONBody struct {
	Active bool    `json:"active"`
	Name   string  `json:"name"`
	Town   *string `json:"town,omitempty"`
}

// CreateServiceProviderJSONBody defines parameters for CreateServiceProvider.
type CreateServiceProviderJSONBody ServiceProvider

// UpdateServiceProviderJSONBody defines parameters for UpdateServiceProvider.
type UpdateServiceProviderJSONBody ServiceProvider

// CreateSessionRequestBody defines body for CreateSession for application/json ContentType.
type CreateSessionJSONRequestBody CreateSessionJSONBody

// ConnectCustomerRequestBody defines body for ConnectCustomer for application/json ContentType.
type ConnectCustomerJSONRequestBody ConnectCustomerJSONBody

// UpdateCustomerRequestBody defines body for UpdateCustomer for application/json ContentType.
type UpdateCustomerJSONRequestBody UpdateCustomerJSONBody

// CreateServiceProviderRequestBody defines body for CreateServiceProvider for application/json ContentType.
type CreateServiceProviderJSONRequestBody CreateServiceProviderJSONBody

// UpdateServiceProviderRequestBody defines body for UpdateServiceProvider for application/json ContentType.
type UpdateServiceProviderJSONRequestBody UpdateServiceProviderJSONBody
