package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/spf13/pflag"
)

const defaultPrefix = "NUTS_"
const defaultDelimiter = "."
const configFileFlag = "configfile"
const defaultConfigFile = "server.config.yaml"
const defaultDBFile = "registry-admin.db"
const defaultHTTPPort = 1303
const defaultNutsNodeAddress = "http://localhost:1323"
const defaultCustomerFile = "customers.json"

func defaultConfig() Config {
	return Config{
		HTTPPort:        defaultHTTPPort,
		DBFile:          defaultDBFile,
		NutsNodeAddress: defaultNutsNodeAddress,
		CustomersFile:   defaultCustomerFile,
	}
}

type Config struct {
	Credentials        Credentials `koanf:"credentials"`
	DBFile             string      `koanf:"dbfile"`
	HTTPPort           int         `koanf:"port"`
	NutsNodeAddress    string      `koanf:"nutsnodeaddr"`
	NutsNodeAPIKeyFile string      `koanf:"nutsnodeapikeyfile"`
	CustomersFile      string      `koanf:"customersfile"`
	Branding           Branding    `koanf:"branding"`
	sessionKey         *ecdsa.PrivateKey
	apiKey             crypto.Signer
}

type Credentials struct {
	Username string `koanf:"username"`
	Password string `koanf:"password" json:"-"` // json omit tag to avoid having it printed in server log
}

type Branding struct {
	// Logo defines a path that points to a file on disk to be used as logo, to be displayed in the application.
	Logo string `koanf:"logo"`
}

func (c Credentials) Empty() bool {
	return len(c.Username) == 0 && len(c.Password) == 0
}

func generateSessionKey() (*ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Printf("failed to generate private key: %s", err)
		return nil, err
	}
	return key, nil
}

func (c Config) Print(writer io.Writer) error {
	if _, err := fmt.Fprintln(writer, "========== CONFIG: =========="); err != nil {
		return err
	}
	var pr Config = c
	data, _ := json.MarshalIndent(pr, "", "  ")
	if _, err := fmt.Println(writer, string(data)); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(writer, "========= END CONFIG ========="); err != nil {
		return err
	}
	return nil
}

func loadConfig() Config {
	flagset := loadFlagSet(os.Args[1:])

	var k = koanf.New(".")

	// Prepare koanf for parsing the config file
	configFilePath := resolveConfigFile(flagset)
	// Check if the file exists
	if _, err := os.Stat(configFilePath); err == nil {
		log.Printf("Loading config from file: %s", configFilePath)
		if err := k.Load(file.Provider(configFilePath), yaml.Parser()); err != nil {
			log.Fatalf("error while loading config from file: %v", err)
		}
	} else {
		log.Printf("Using default config because no file was found at: %s", configFilePath)
	}

	// load env flags, can't return error
	_ = k.Load(envProvider(), nil)

	config := defaultConfig()
	var err error
	sessionKey, err := generateSessionKey()
	if err != nil {
		log.Fatalf("unable to generate session key: %v", err)
	}
	config.sessionKey = sessionKey

	// Unmarshal values of the config file into the config struct, potentially replacing default values
	if err := k.Unmarshal("", &config); err != nil {
		log.Fatalf("error while unmarshalling config: %v", err)
	}

	// Load the API key
	if len(config.NutsNodeAPIKeyFile) > 0 {
		bytes, err := os.ReadFile(config.NutsNodeAPIKeyFile)
		if err != nil {
			log.Fatalf("error while reading private key file: %v", err)
		}
		config.apiKey, err = pemToPrivateKey(bytes)
		if err != nil {
			log.Fatalf("error while decoding private key file: %v", err)
		}
	}

	return config
}

func loadFlagSet(args []string) *pflag.FlagSet {
	f := pflag.NewFlagSet("config", pflag.ContinueOnError)
	f.String(configFileFlag, defaultConfigFile, "Nuts config file")
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}
	f.Parse(args)
	return f
}

// resolveConfigFile resolves the path of the config file using the following sources:
// 1. commandline params (using the given flags)
// 2. environment vars,
// 3. default location.
func resolveConfigFile(flagset *pflag.FlagSet) string {

	k := koanf.New(defaultDelimiter)

	// load env flags, can't return error
	_ = k.Load(envProvider(), nil)

	// load cmd flags, without a parser, no error can be returned
	_ = k.Load(posflag.Provider(flagset, defaultDelimiter, k), nil)

	configFile := k.String(configFileFlag)
	return configFile
}

func envProvider() *env.Env {
	return env.Provider(defaultPrefix, defaultDelimiter, func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, defaultPrefix)), "_", defaultDelimiter, -1)
	})
}

// pemToPrivateKey converts a PEM encoded private key to a Signer interface. It supports EC, RSA and PKIX PEM encoded strings
func pemToPrivateKey(bytes []byte) (signer crypto.Signer, err error) {
	key, _ := ssh.ParseRawPrivateKey(bytes)
	if key == nil {
		err = errors.New("failed to decode PEM file")
		return
	}

	switch k := key.(type) {
	case *rsa.PrivateKey:
		signer = k
	case *ecdsa.PrivateKey:
		signer = k
	default:
		err = fmt.Errorf("unknown private key type: %s", k)
	}

	return
}
