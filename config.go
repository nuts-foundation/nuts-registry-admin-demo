package main

import (
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
const defaultConfigFile = "default.config.yaml"

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
	flagset := loadFlagSet(os.Args[1:])

	var k = koanf.New(".")

	if err := k.Load(file.Provider(resolveConfigFile(flagset)), yaml.Parser()); err != nil {
		log.Fatalf("error while loading config from file: %v", err)
	}
	config := defaultConfig()
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
