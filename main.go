package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

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


func main() {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Check if we use live mode from the file system or using embedded files
	useFS := len(os.Args) > 1 && os.Args[1] == "live"

	assetHandler := http.FileServer(getFileSystem(useFS))
	e.GET("/*", echo.WrapHandler(assetHandler))
	e.GET("/api/customers", func(ctx echo.Context) error {
		customers := []map[string]string{
			{
        "identifier":"notenboom.vendor.com",
        "name":"Zorginstelling de notenboom",
        "town":"Amsterdam",
        "did":"did:nuts:123",
        "active":"1",
      },
			{
        "identifier":"nootjes.vendor.com",
        "name":"Verpleeghuis de nootjes",
        "town":"Maastricht",
        "did":"did:nuts:456",
        "active":"",
      },
		}
		return ctx.JSON(200, customers)
	})

	e.Logger.Fatal(e.Start(":1303"))
}
