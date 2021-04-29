package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/lestrrat-go/jwx/jwt/openid"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain/customers"
)

type Wrapper struct {
	Auth auth
	SPRepo domain.ServiceProviderRepository
	CustomerService customers.Service
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
	credentials := CreateSessionRequest{}
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

	return ctx.JSON(200, CreateSessionResponse{Token: string(token)})
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
		return err
	}
	response := CustomersResponse{}
	for _, c := range customers {
		customer := Customer{
			Did:  c.DID,
			Id:   c.ID,
			Name: c.Name,
		}
		response = append(response, customer)
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
	response := ServiceProvider{
		Email: sp.Email,
		Id:    sp.ID,
		Name:  sp.Name,
		Phone: sp.Phone,
	}
	return ctx.JSON(200, response)
}

func (w Wrapper) CreateServiceProvider(ctx echo.Context) error {
	_, err := w.checkAuthorization(ctx)
	if err != nil {
		return err
	}
	spRequest := ServiceProvider{}
	if err := ctx.Bind(&spRequest); err != nil {
		return err
	}
	sp := domain.ServiceProvider{
		ID:    spRequest.Id,
		Name:  spRequest.Name,
		Email: spRequest.Email,
		Phone: spRequest.Phone,
	}
	if err := w.SPRepo.CreateOrUpdate(sp); err != nil {
		return err
	}
	return ctx.JSON(201, spRequest)
}

func (w Wrapper) UpdateServiceProvider(ctx echo.Context) error {
	_, err := w.checkAuthorization(ctx)
	if err != nil {
		return err
	}
	spRequest := ServiceProvider{}
	if err := ctx.Bind(&spRequest); err != nil {
		return err
	}
	sp := domain.ServiceProvider{
		ID:    spRequest.Id,
		Name:  spRequest.Name,
		Email: spRequest.Email,
		Phone: spRequest.Phone,
	}
	if err := w.SPRepo.CreateOrUpdate(sp); err != nil {
		return err
	}
	return ctx.JSON(200, spRequest)
}

func (w Wrapper) ConnectCustomer(ctx echo.Context) error {
	token, err := w.checkAuthorization(ctx)
	if err != nil {
		return err
	}

	connectReq := ConnectCustomerRequest{}
	if err := ctx.Bind(&connectReq); err != nil {
		return err
	}

	if user, ok := token.Get(openid.EmailKey); ok {
		ctx.Logger().Printf("Customer with id %s connected by: %s", user, connectReq.Id)
	} else {
		return echo.NewHTTPError(http.StatusForbidden, "unknown user")
	}

	customer, err := w.CustomerService.ConnectCustomer(connectReq.Id, connectReq.Name)
	if err != nil {
		return err
	}
	return ctx.JSON(200, customer)
}
