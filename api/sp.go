package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	ssi "github.com/nuts-foundation/go-did"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain"
)

func (w Wrapper) GetServiceProvider(ctx echo.Context) error {
	serviceProvider, err := w.SPService.Get()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if serviceProvider == nil {
		return ctx.NoContent(http.StatusNotFound)
	}
	return ctx.JSON(http.StatusOK, serviceProvider)
}

func (w Wrapper) UpdateServiceProvider(ctx echo.Context) error {
	serviceProvider := domain.ServiceProvider{}
	if err := ctx.Bind(&serviceProvider); err != nil {
		return err
	}
	res, err := w.SPService.CreateOrUpdate(serviceProvider)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	// Make sure NutsComm service is registered on customers' DID documents
	customers, err := w.CustomerService.Repository.All()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	for _, customer := range customers {
		if customer.Did == nil {
			continue
		}
		err := w.CustomerService.RegisterNutsCommService(customer.Id, serviceProvider.Id)
		if err != nil {
			log.Printf("Couldn't register NutsComm endpoint on customer DID (did=%s): %v", *customer.Did, err)
		}
	}
	return ctx.JSON(http.StatusOK, res)
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
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if serviceProvider == nil {
		return ctx.NoContent(http.StatusNotFound)
	}
	endpoints, err := w.SPService.GetEndpoints(*serviceProvider)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, endpoints)
}

func (w Wrapper) GetServices(ctx echo.Context) error {
	services, err := w.SPService.GetServices()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, services)
}

func (w Wrapper) AddService(ctx echo.Context) error {
	service := domain.ServiceProperties{}
	ctx.Bind(&service)
	addedService, err := w.SPService.AddService(service)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, addedService)
}
