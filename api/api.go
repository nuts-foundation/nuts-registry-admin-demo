package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwt/openid"
)

type Wrapper struct {
	Auth auth
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
	authHeader := ctx.Request().Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return echo.NewHTTPError(http.StatusForbidden, "Authorization header must contain 'Bearer <token>'")
	}
	tokenStr := strings.Split(authHeader, " ")[1]
	token, err := w.Auth.ValidateJWT([]byte(tokenStr))
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("invalid token: %s", err))
	}
	if user, ok := token.Get(openid.EmailKey); ok {
		ctx.Logger().Print("Customers requested by: %s", user)
	} else {
		return echo.NewHTTPError(http.StatusForbidden, "unknown user")
	}

	customers := []map[string]string{
		{"name": "Zorginstelling de notenboom", "did": "did:nuts:123"},
		{"name": "Verpleehuis de nootjes", "did": "did:nuts:456"},
	}
	return ctx.JSON(200, customers)
}
