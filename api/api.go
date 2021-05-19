package api

import (
	"fmt"
	"net/http"

	ssi "github.com/nuts-foundation/go-did"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain/sp"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/go-did/did"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain/credentials"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain/customers"
)

type Wrapper struct {
	Auth              auth
	SPService         sp.Service
	CustomerService   customers.Service
	CredentialService credentials.Service
}

func (w Wrapper) CreateSession(ctx echo.Context) error {
	sessionRequest := domain.CreateSessionRequest{}
	if err := ctx.Bind(&sessionRequest); err != nil {
		return err
	}

	if !w.Auth.CheckCredentials(sessionRequest.Username, sessionRequest.Password) {
		return echo.NewHTTPError(http.StatusForbidden, "invalid sessionRequest")
	}

	token, err := w.Auth.CreateJWT(sessionRequest.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(200, domain.CreateSessionResponse{Token: string(token)})
}

func (w Wrapper) GetCustomers(ctx echo.Context) error {
	allCustomers, err := w.CustomerService.Repository.All()
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	for i, c := range allCustomers {
		credentialsForCustomer, err := w.CredentialService.GetCredentials(c)
		if err != nil {
			return echo.NewHTTPError(500, err.Error())
		}
		allCustomers[i].Active = len(credentialsForCustomer) > 0
	}

	response := domain.CustomersResponse{}
	for _, c := range allCustomers {
		response = append(response, c)
	}
	return ctx.JSON(200, response)
}

func (w Wrapper) GetServiceProvider(ctx echo.Context) error {
	serviceProvider, err := w.SPService.Get()
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	if serviceProvider == nil {
		return ctx.NoContent(404)
	}
	return ctx.JSON(200, serviceProvider)
}

func (w Wrapper) UpdateServiceProvider(ctx echo.Context) error {
	serviceProvider := domain.ServiceProvider{}
	if err := ctx.Bind(&serviceProvider); err != nil {
		return err
	}
	res, err := w.SPService.CreateOrUpdate(serviceProvider)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return ctx.JSON(200, res)
}

func (w Wrapper) RegisterEndpoint(ctx echo.Context) error {
	ep := domain.Endpoint{}
	if err := ctx.Bind(&ep); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	err := w.SPService.RegisterEndpoint(ep)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusCreated)
}

func (w Wrapper) ConnectCustomer(ctx echo.Context) error {
	connectReq := domain.ConnectCustomerRequest{}
	if err := ctx.Bind(&connectReq); err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	if len(connectReq.Id) == 0 || len(connectReq.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "name and id must be provided")
	}

	town := ""
	if connectReq.Town != nil {
		town = *connectReq.Town
	}

	serviceProvider, err := w.SPService.Get()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("unable to fetch service provider ID: %w", err))
	}
	spID, err := did.ParseDID(serviceProvider.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("service provider not correctly configured: DID is invalid: %w", err))
	}

	customer, err := w.CustomerService.ConnectCustomer(connectReq.Id, connectReq.Name, town, *spID)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return ctx.JSON(200, customer)
}

func (w Wrapper) UpdateCustomer(ctx echo.Context, id string) error {
	req := struct {
		Name   string
		Active bool
		Town   string
	}{}
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	if len(req.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "name")
	}

	customer, err := w.CustomerService.Repository.Update(id, func(c domain.Customer) (*domain.Customer, error) {
		c.Name = req.Name
		if len(req.Town) >= 0 {
			c.Town = &req.Town
		}
		if err := w.CredentialService.ManageNutsOrgCredential(c, req.Active); err != nil {
			return nil, err
		}
		return &c, nil
	})
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return ctx.JSON(200, customer)
}

func (w Wrapper) GetCustomer(ctx echo.Context, id string) error {
	customer, err := w.CustomerService.Repository.FindByID(id)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	if customer == nil {
		return ctx.NoContent(404)
	}

	credentialsForCustomer, err := w.CredentialService.GetCredentials(*customer)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	customer.Active = len(credentialsForCustomer) > 0
	return ctx.JSON(200, customer)
}

func (w Wrapper) GetCredentialIssuers(ctx echo.Context) error {
	res, err := w.CredentialService.GetCredentialIssuers([]string{"NutsOrganizationCredential"})
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return ctx.JSON(200, res)
}

func (w Wrapper) UpdateCredentialIssuer(ctx echo.Context, CredentialType string, didStr string) error {
	var request = struct {
		Trusted bool
	}{}
	ctx.Bind(&request)
	id, err := ssi.ParseURI(didStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	issuerTrust, err := w.CredentialService.ManageIssuerTrust(CredentialType, *id, request.Trusted)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return ctx.JSON(200, issuerTrust)
}

