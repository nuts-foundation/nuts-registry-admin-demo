package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
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

func defaultConfig() Config {
	return Config{
		HTTPPort: defaultHTTPPort,
		DBFile: defaultDBFile,
	}
}

type Config struct {
	Credentials struct {
		Username string `koanf:"username"`
		Password string `koanf:"password"`
	}
	DBFile     string `koanf:"dbfile"`
	HTTPPort   int    `koanf:"port"`
	SessionKey *ecdsa.PrivateKey
}

func generateSessionKey() (*ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Printf("failed to generate private key: %s", err)
		return nil, err
	}
	return key, nil
}

func loadConfig() Config {
	flagset := loadFlagSet(os.Args[1:])

	var k = koanf.New(".")

	if err := k.Load(file.Provider(resolveConfigFile(flagset)), yaml.Parser()); err != nil {
		log.Fatalf("error while loading config from file: %v", err)
	}

	config := defaultConfig()
	var err error
	sessionKey, err := generateSessionKey()
	if err != nil {
		log.Fatalf("unable to generate session key: %v", err)
	}
	config.SessionKey = sessionKey

	if err := k.Unmarshal("", &config); err != nil {
		log.Fatalf("error while unmarshalling config: %v", err)
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

	// load env flags
	e := env.Provider(defaultPrefix, defaultDelimiter, func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, defaultPrefix)), "_", defaultDelimiter, -1)
	})
	// can't return error
	_ = k.Load(e, nil)

	// load cmd flags, without a parser, no error can be returned
	_ = k.Load(posflag.Provider(flagset, defaultDelimiter, k), nil)

	configFile := k.String(configFileFlag)
	log.Printf("using config: %s", configFile)
	return configFile
}
