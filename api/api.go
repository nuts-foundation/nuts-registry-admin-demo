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

func (w Wrapper) CheckSession(ctx echo.Context) error {
	// If this function is reached, it means the session is still valid
	return ctx.NoContent(http.StatusNoContent)
}

func (w Wrapper) CreateSession(ctx echo.Context) error {
	sessionRequest := domain.CreateSessionRequest{}
	if err := ctx.Bind(&sessionRequest); err != nil {
		return err
	}

	if !w.Auth.CheckCredentials(sessionRequest.Username, sessionRequest.Password) {
		return echo.NewHTTPError(http.StatusForbidden, "invalid credentials")
	}

	token, err := w.Auth.CreateJWT(sessionRequest.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, domain.CreateSessionResponse{Token: string(token)})
}

func (w Wrapper) GetCustomers(ctx echo.Context) error {
	allCustomers, err := w.CustomerService.Repository.All()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for i, c := range allCustomers {
		credentialsForCustomer, err := w.CredentialService.GetCredentials(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		allCustomers[i].Active = len(credentialsForCustomer) > 0
	}

	response := domain.CustomersResponse{}
	for _, c := range allCustomers {
		response = append(response, c)
	}
	return ctx.JSON(http.StatusOK, response)
}

func (w Wrapper) ConnectCustomer(ctx echo.Context) error {
	customer := &domain.Customer{}
	if err := ctx.Bind(customer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(customer.Id) == 0 || len(customer.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "name and id must be provided")
	}

	serviceProvider, err := w.SPService.Get()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("unable to fetch service provider ID: %w", err))
	}
	spID, err := did.ParseDID(serviceProvider.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("service provider not correctly configured: DID is invalid: %w", err))
	}

	customer, err = w.CustomerService.ConnectCustomer(*customer, *spID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, customer)
}

func (w Wrapper) UpdateCustomer(ctx echo.Context, id string) error {
	req := domain.Customer{}
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	if len(req.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "name")
	}

	customer, err := w.CustomerService.Repository.Update(id, func(c domain.Customer) (*domain.Customer, error) {
		c.Name = req.Name
		c.City = req.City
		c.Domain = req.Domain
		if err := w.CredentialService.ManageNutsOrgCredential(c, req.Active); err != nil {
			return nil, err
		}
		return &c, nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, customer)
}

func (w Wrapper) GetCustomer(ctx echo.Context, id string) error {
	customer, err := w.CustomerService.Repository.FindByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if customer == nil {
		return ctx.NoContent(http.StatusNotFound)
	}

	credentialsForCustomer, err := w.CredentialService.GetCredentials(*customer)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	customer.Active = len(credentialsForCustomer) > 0
	return ctx.JSON(http.StatusOK, customer)
}

func (w Wrapper) GetCredentialIssuers(ctx echo.Context) error {
	res, err := w.CredentialService.GetCredentialIssuers([]string{"NutsOrganizationCredential"})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, res)
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
	return ctx.JSON(http.StatusOK, issuerTrust)
}

func (w Wrapper) SearchOrganizations(ctx echo.Context) error {
	params := domain.SearchOrganizationsJSONBody{}
	if err := ctx.Bind(&params); err != nil {
		return err
	}
	result, err := w.CredentialService.SearchOrganizations(params.Name, params.City)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, result)
}

func (w Wrapper) GetServicesForCustomer(ctx echo.Context, customerID string) error {
	services, err := w.CustomerService.GetServices(customerID)
	if err != nil {
		return err
	}
	// make sure the response is always initialized to ensure [] instead of null json
	response := make([]did.Service, len(services))
	for i, s := range services {
		response[i] = s
	}
	return ctx.JSON(http.StatusOK, response)
}

func (w Wrapper) EnableCustomerService(ctx echo.Context, customerID string) error {
	req := domain.EnableCustomerServiceJSONRequestBody{}
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	if err := w.CustomerService.EnableService(customerID, req.Did, req.Type); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (w Wrapper) DisableCustomerService(ctx echo.Context, customerID string, serviceType string) error {
	if err := w.CustomerService.DisableService(customerID, serviceType); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}
