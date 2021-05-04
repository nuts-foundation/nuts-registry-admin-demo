package main

import (
	"crypto/elliptic"
	"crypto/sha1"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nuts-foundation/nuts-registry-admin-demo/api"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain/credentials"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain/customers"
	bolt "go.etcd.io/bbolt"
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
	config := loadConfig()
	config.Print(log.Writer())
	// load bbolt db
	db, err := bolt.Open(config.DBFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())

	// Initialize Auth
	var account api.UserAccount
	if config.Credentials.Empty() {
		account = generateDefaultAccount(config)
		log.Printf("Authentication credentials not configured, so they were generated (user=%s, password=%s)", account.Username, account.Password)
	} else {
		account = api.UserAccount{Username: config.Credentials.Username, Password: config.Credentials.Password}
	}
	auth := api.NewAuth(config.sessionKey, []api.UserAccount{account})

	// Initialize repos
	spRepo := domain.ServiceProviderRepository{NodeAddr: config.NutsNodeAddress, DB: db}
	cRepo := customers.NewDB(config.CustomersFile)

	// Initialize services
	customerService := customers.Service{
		NutsNodeAddr: config.NutsNodeAddress,
		Repository:   cRepo,
	}
	credentialService := credentials.Service{
		NutsNodeAddr: config.NutsNodeAddress,
		SPRepository: spRepo,
	}

	// Initialize wrapper
	apiWrapper := api.Wrapper{Auth: auth, SPRepo: spRepo, CustomerService: customerService, CredentialService: credentialService}

	api.RegisterHandlers(e, apiWrapper)

	// Setup asset serving:
	// Check if we use live mode from the file system or using embedded files
	useFS := len(os.Args) > 1 && os.Args[1] == "live"
	assetHandler := http.FileServer(getFileSystem(useFS))
	e.GET("/*", echo.WrapHandler(assetHandler))

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.HTTPPort)))
}

func generateDefaultAccount(config Config) api.UserAccount {
	pkHashBytes := sha1.Sum(elliptic.Marshal(config.sessionKey.Curve, config.sessionKey.X, config.sessionKey.Y))
	return api.UserAccount{Username: "demo@nuts.nl", Password: hex.EncodeToString(pkHashBytes[:])}
}
