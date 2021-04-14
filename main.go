package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func main() {
	config := loadConfig()

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
		if credentials.Username != config.Credentials.Username|| credentials.Password != config.Credentials.Password {
			return echo.NewHTTPError(http.StatusForbidden, "invalid credentials")
		}
		return ctx.JSON(200, map[string]string{"token": "123"})
	})
	e.GET("/api/customers", func(ctx echo.Context) error {
		customers := []map[string]string{
			{"name": "Zorginstelling de notenboom", "did": "did:nuts:123"},
			{"name": "Verpleehuis de nootjes", "did": "did:nuts:456"},
		}
		return ctx.JSON(200, customers)
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.HTTPPort)))
}
