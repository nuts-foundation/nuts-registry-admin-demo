package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/lestrrat-go/jwx/jwt/openid"
)

const assetPath = "web/dist"

//go:embed web/dist/*
var embeddedFiles embed.FS

func getFileSystem(useFS bool) http.FileSystem {
	if useFS {
		log.Print("using live mode")
		return http.FS(os.DirFS(assetPath))
	}

	log.Print("using embed mode")
	fsys, err := fs.Sub(embeddedFiles, assetPath)
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

func defaultConfig() Config {
	return Config{
		HTTPPort: 1303,
	}
}

type Config struct {
	Credentials struct {
		Username string `koanf:"username"`
		Password string `koanf:"password"`
	}
	HTTPPort int `koanf:"port"`
}

func loadConfig() Config {
	var k = koanf.New(".")

	if err := k.Load(file.Provider("server.config.yaml"), yaml.Parser()); err != nil {
		log.Fatalf("error while loading config from file: %v", err)
	}
	config := defaultConfig()
	if err := k.Unmarshal("", &config); err != nil {
		log.Fatalf("error while unmarshalling config: %v", err)
	}
	return config
}

func createJWT(email string) ([]byte, error) {
	t := openid.New()
	t.Set(jwt.IssuedAtKey, time.Now())
	// session is valid for 20 minutes
	t.Set(jwt.ExpirationKey, time.Now().Add(20*time.Minute))
	t.Set(openid.EmailKey, email)

	signed, err := jwt.Sign(t, jwa.ES256, sessionKey)
	if err != nil {
		log.Printf("failed to sign token: %s", err)
		return nil, err
	}
	return signed, nil
}

func validateJWT(token []byte) (jwt.Token, error) {
	pubKey := sessionKey.(*ecdsa.PrivateKey).PublicKey
	t, err := jwt.Parse(token, jwt.WithVerify(jwa.ES256, pubKey), jwt.WithValidate(true))
	if err != nil {
		log.Printf("unable to parse token: %s", err)
		return nil, err
	}
	return t, nil
}

var sessionKey crypto.PrivateKey

func generateSessionKey() (crypto.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Printf("failed to generate private key: %s", err)
		return nil, err
	}
	return key, nil
}

func main() {
	config := loadConfig()
	var err error
	sessionKey, err = generateSessionKey()
	if err != nil {
		log.Fatalf("unable to generate session key: %v", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Check if we use live mode from the file system or using embedded files
	useFS := len(os.Args) > 1 && os.Args[1] == "live"

	assetHandler := http.FileServer(getFileSystem(useFS))
	e.GET("/*", echo.WrapHandler(assetHandler))
	e.POST("/api/auth", func(ctx echo.Context) error {
		credentials := struct {
			Username string
			Password string
		}{}
		err := ctx.Bind(&credentials)
		if err != nil {
			return err
		}
		if credentials.Username != config.Credentials.Username || credentials.Password != config.Credentials.Password {
			return echo.NewHTTPError(http.StatusForbidden, "invalid credentials")
		}
		token, err := createJWT(credentials.Username)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return ctx.JSON(200, map[string]string{"token": string(token)})
	})
	e.GET("/api/customers", func(ctx echo.Context) error {
		authHeader := ctx.Request().Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return echo.NewHTTPError(http.StatusForbidden, "Authorization header must contain 'Bearer <token>'")
		}
		tokenStr := strings.Split(authHeader, " ")[1]
		token, err := validateJWT([]byte(tokenStr))
		if err != nil {
			return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("invalid token: %s", err))
		}
		if user, ok := token.Get(openid.EmailKey); ok {
			log.Printf("Customers requested by: %s", user)
		} else {
			return echo.NewHTTPError(http.StatusForbidden, "unknown user")
		}

		customers := []map[string]string{
			{"name": "Zorginstelling de notenboom", "did": "did:nuts:123"},
			{"name": "Verpleehuis de nootjes", "did": "did:nuts:456"},
		}
		return ctx.JSON(200, customers)
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.HTTPPort)))
}
