package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/lestrrat-go/jwx/jwt/openid"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain/credentials"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain/customers"
)

type Wrapper struct {
	Auth              auth
	SPRepo            domain.ServiceProviderRepository
	CustomerService   customers.Service
	CredentialService credentials.Service
}

func (w Wrapper) checkAuthorization(ctx echo.Context) (jwt.Token, error) {
	authHeader := ctx.Request().Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, echo.NewHTTPError(http.StatusForbidden, "Authorization header must contain 'Bearer <token>'")
	}
	tokenStr := strings.Split(authHeader, " ")[1]
	token, err := w.Auth.ValidateJWT([]byte(tokenStr))
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("invalid token: %s", err))
	}
	return token, nil
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
	token, err := w.checkAuthorization(ctx)
	if err != nil {
		return err
	}
	if user, ok := token.Get(openid.EmailKey); ok {
		ctx.Logger().Printf("Customers requested by: %s", user)
	} else {
		return echo.NewHTTPError(http.StatusForbidden, "unknown user")
	}

	customers, err := w.CustomerService.Repository.All()
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	for _, c := range customers {
		credentials, err := w.CredentialService.GetCredentials(c)
		if err != nil {
			return echo.NewHTTPError(500, err.Error())
		}
		if len(credentials) > 0 {
			c.Active = true
		}
	}

	response := domain.CustomersResponse{}
	for _, c := range customers {
		response = append(response, c)
	}
	return ctx.JSON(200, response)
}

func (w Wrapper) GetServiceProvider(ctx echo.Context) error {
	_, err := w.checkAuthorization(ctx)
	if err != nil {
		return err
	}
	sp, err := w.SPRepo.Get()
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	if sp == nil {
		return ctx.NoContent(404)
	}
	return ctx.JSON(200, sp)
}

func (w Wrapper) CreateServiceProvider(ctx echo.Context) error {
	_, err := w.checkAuthorization(ctx)
	if err != nil {
		return err
	}
	sp := domain.ServiceProvider{}
	if err := ctx.Bind(&sp); err != nil {
		return err
	}
	if err := w.SPRepo.CreateOrUpdate(sp); err != nil {
		return err
	}
	return ctx.JSON(201, sp)
}

func (w Wrapper) UpdateServiceProvider(ctx echo.Context) error {
	_, err := w.checkAuthorization(ctx)
	if err != nil {
		return err
	}
	sp := domain.ServiceProvider{}
	if err := ctx.Bind(&sp); err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	if err := w.SPRepo.CreateOrUpdate(sp); err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return ctx.JSON(200, sp)
}

func (w Wrapper) ConnectCustomer(ctx echo.Context) error {
	token, err := w.checkAuthorization(ctx)
	if err != nil {
		return err
	}

	connectReq := domain.ConnectCustomerRequest{}
	if err := ctx.Bind(&connectReq); err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	if user, ok := token.Get(openid.EmailKey); ok {
		ctx.Logger().Printf("Customer with id %s connected by: %s", user, connectReq.Id)
	} else {
		return echo.NewHTTPError(http.StatusForbidden, "unknown user")
	}

	if len(connectReq.Id) == 0 || len(connectReq.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "name and id must be provided")
	}

	customer, err := w.CustomerService.ConnectCustomer(connectReq.Id, connectReq.Name)
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
		if len(req.Town) > 0 {
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
	return ctx.JSON(200, customer)
}
