package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nuts-foundation/nuts-registry-admin-demo/api"
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

func generateSessionKey() (*ecdsa.PrivateKey, error) {
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
	sessionKey, err := generateSessionKey()
	if err != nil {
		log.Fatalf("unable to generate session key: %v", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	auth := api.NewAuth(sessionKey, []api.UserAccount{{Username: config.Credentials.Username, Password: config.Credentials.Password}})
	apiWrapper := api.Wrapper{Auth: auth}

	api.RegisterHandlers(e, apiWrapper)

	// Check if we use live mode from the file system or using embedded files
	useFS := len(os.Args) > 1 && os.Args[1] == "live"

	assetHandler := http.FileServer(getFileSystem(useFS))
	e.GET("/*", echo.WrapHandler(assetHandler))
	//e.GET("/api/customers", func(ctx echo.Context) error {
	//})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.HTTPPort)))
}
