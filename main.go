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
	"strings"
	"time"

	didmanAPI "github.com/nuts-foundation/nuts-node/didman/api/v1"
	vdrAPI "github.com/nuts-foundation/nuts-node/vdr/api/v1"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain/sp"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/nuts-foundation/nuts-registry-admin-demo/api"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain/credentials"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain/customers"
	bolt "go.etcd.io/bbolt"
)

const assetPath = "web/dist"

//go:embed web/dist/*
var embeddedFiles embed.FS

const apiTimeout = 10 * time.Second

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
	loggerConfig := middleware.DefaultLoggerConfig
	loggerConfig.Skipper = requestsStatusEndpoint
	e.Use(middleware.LoggerWithConfig(loggerConfig))
	//e.Debug = true
	//e.Use(middleware.Recover())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper: func(c echo.Context) bool {
			protectedPaths := []string{
				"/web/private",
			}
			for _, path := range protectedPaths {
				if strings.HasPrefix(c.Request().RequestURI, path) {
					return false
				}
			}
			return true
		},
		SigningKey:    &config.sessionKey.PublicKey,
		SigningMethod: jwa.ES256.String(),
	}))
	e.HTTPErrorHandler = httpErrorHandler

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
	vdrClient := vdrAPI.HTTPClient{
		ServerAddress: config.NutsNodeAddress,
		Timeout:       apiTimeout,
	}
	didmanClient := didmanAPI.HTTPClient{
		ServerAddress: config.NutsNodeAddress,
		Timeout:       apiTimeout,
	}
	spService := sp.Service{
		Repository:   sp.NewBBoltRepository(db),
		VDRClient:    vdrClient,
		DIDManClient: didmanClient,
	}

	// Initialize services
	customerService := customers.Service{
		VDRClient:    vdrClient,
		Repository:   customers.NewFlatFileRepository(config.CustomersFile),
		DIDManClient: didmanClient,
	}
	credentialService := credentials.Service{
		NutsNodeAddr: config.NutsNodeAddress,
		SPService:    spService,
		DIDManClient: didmanClient,
	}

	// Initialize wrapper
	apiWrapper := api.Wrapper{Auth: auth, SPService: spService, CustomerService: customerService, CredentialService: credentialService}

	api.RegisterHandlers(e, apiWrapper)

	// Setup asset serving:
	// Check if we use live mode from the file system or using embedded files
	useFS := len(os.Args) > 1 && os.Args[1] == "live"
	assetHandler := http.FileServer(getFileSystem(useFS))
	e.GET("/branding/logo", (&api.LogoHandler{FilePath: config.Branding.Logo}).Handle)
	e.GET("/status", func(context echo.Context) error {
		return context.String(http.StatusOK, "OK")
	})
	e.GET("/*", echo.WrapHandler(assetHandler))

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.HTTPPort)))
}

func generateDefaultAccount(config Config) api.UserAccount {
	pkHashBytes := sha1.Sum(elliptic.Marshal(config.sessionKey.Curve, config.sessionKey.X, config.sessionKey.Y))
	return api.UserAccount{Username: "demo@nuts.nl", Password: hex.EncodeToString(pkHashBytes[:])}
}

// httpErrorHandler includes the err.Err() string in a { "error": "msg" } json hash
func httpErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)
	type Map map[string]interface{}

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
		if he.Internal != nil {
			err = fmt.Errorf("%v, %v", err, he.Internal)
		}
	} else {
		msg = err.Error()
	}

	if _, ok := msg.(string); ok {
		msg = Map{"error": msg}
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, msg)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}

func requestsStatusEndpoint(context echo.Context) bool {
	return context.Request().RequestURI == "/status"
}
