package api

import (
	"fmt"
	"net/http"

	"github.com/nuts-foundation/nuts-registry-admin-demo/domain/sp"

	"github.com/labstack/echo/v4"
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
	credentials := domain.CreateSessionRequest{}
	if err := ctx.Bind(&credentials); err != nil {
		return err
	}

	if !w.Auth.CheckCredentials(credentials.Username, credentials.Password) {
		return echo.NewHTTPError(http.StatusForbidden, "invalid credentials")
	}

	token, err := w.Auth.CreateJWT(credentials.Username)
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
	sp, err := w.SPService.Get()
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	if sp == nil {
		return ctx.NoContent(404)
	}
	return ctx.JSON(200, sp)
}

func (w Wrapper) CreateServiceProvider(ctx echo.Context) error {
	sp := domain.ServiceProvider{}
	if err := ctx.Bind(&sp); err != nil {
		return err
	}
	res, err := w.SPService.CreateOrUpdate(sp)
	if err != nil {
		return err
	}
	return ctx.JSON(201, res)
}

func (w Wrapper) UpdateServiceProvider(ctx echo.Context) error {
	sp := domain.ServiceProvider{}
	if err := ctx.Bind(&sp); err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	res, err := w.SPService.CreateOrUpdate(sp)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return ctx.JSON(200, res)
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

	spID, err := w.SPService.Repository.Get()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("unable to fetch service provider ID: %w", err))
	}

	customer, err := w.CustomerService.ConnectCustomer(connectReq.Id, connectReq.Name, town, spID)
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
	ctx.Bind(&req)

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
		ctx.JSON(500, err.Error())
	}
	return ctx.JSON(200, customer)
}

func (w Wrapper) GetCustomer(ctx echo.Context, id string) error {
	customer, err := w.CustomerService.Repository.FindByID(id)
	if err != nil {
		ctx.JSON(500, err.Error())
	}
	if customer == nil {
		ctx.NoContent(404)
	}

	credentialsForCustomer, err := w.CredentialService.GetCredentials(*customer)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	customer.Active = len(credentialsForCustomer) > 0
	return ctx.JSON(200, customer)
}
