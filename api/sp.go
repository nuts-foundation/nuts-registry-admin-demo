package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	ssi "github.com/nuts-foundation/go-did"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain"
)

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
	ep := domain.EndpointProperties{}
	if err := ctx.Bind(&ep); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	err := w.SPService.RegisterEndpoint(ep)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusCreated)
}

func (w Wrapper) DeleteEndpoint(ctx echo.Context, idStr string) error {
	id, err := ssi.ParseURI(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("invalid endpoint ID: %w", err))
	}

	if err := w.SPService.DeleteEndpoint(*id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusNoContent)
}
func (w Wrapper) GetEndpoints(ctx echo.Context) error {
	serviceProvider, err := w.SPService.Get()
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	if serviceProvider == nil {
		return ctx.NoContent(404)
	}
	endpoints, err := w.SPService.Endpoints(*serviceProvider)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return ctx.JSON(200, endpoints)
}
