// Package config ensures the app is configured properly
package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

var config *Config

// Config is the CLI options wrapped in a struct
type Config struct {
	ConfigFile string

	Debug         bool   `yaml:"debug"`
	Listen        string `yaml:"listen"`
	Port          uint   `yaml:"port"`
	MimirEndpoint string `yaml:"mimir_endpoint"`

	zapConfig *zap.Config
	zapLogger *zap.Logger
	zapSugar  *zap.SugaredLogger
}

// NewConfig returns a Config struct that can be used to reference configuration
// NewConfig does the following:
//   - Runs initCLI (sets and read CLI switches)
//   - Runs initConfig (reads config from files)
func NewConfig(name string) (*Config, error) {
	if config == nil {
		config = &Config{}

		// Stage 1: Setup CLI flags
		initCLI(config)

		// Stage 2: Setup config file
		if err := initConfig(config); err != nil {
			return nil, err
		}

		// Stage 3: initialize logging
		initLogging(config)

		// Stage 4: Setup env vars
		if err := initEnv(config); err != nil {
			return nil, err
		}

		// Stage 5: Check configuration
		if err := checkConfig(config); err != nil {
			return nil, err
		}
	}

	return config, nil
}

// initCLI initializes CLI switches
func initCLI(config *Config) {
	flag.BoolVar(&config.Debug, "debug", true, "Debug output")
	flag.StringVar(&config.Listen, "listen", "0.0.0.0", "Listen address")
	flag.UintVar(&config.Port, "port", 8080, "Listen port") //nolint
	flag.StringVar(&config.ConfigFile, "config", "./config.yaml", "Path for the config.yaml configuration file")
	flag.StringVar(&config.MimirEndpoint, "remote-write-endpoint", "http://localhost:8080", "Endpoint to remotely write metrics to")
	flag.Parse()
}

// initConfig initializes file config and converges CLI, file, and env var
func initConfig(config *Config) error {
	// Config file: config.yaml
	data, err := os.ReadFile(config.ConfigFile)
	if err != nil {
		return err
	}

	if err1 := yaml.Unmarshal(data, config); err1 != nil {
		return err
	}

	return nil
}

func initLogging(config *Config) {
	config.zapConfig = &zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:         "console",
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			FunctionKey:    "F",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	if config.Debug {
		config.zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	config.zapLogger = zap.Must(config.zapConfig.Build())
	config.zapSugar = config.zapLogger.Sugar()
	zap.ReplaceGlobals(config.zapLogger)
}

func initEnv(config *Config) error {
	listen := os.Getenv("LISTEN")
	port := os.Getenv("PORT")
	remoteWriteEndpoint := os.Getenv("REMOTE_WRITE_ENDPOINT")

	if listen != "" {
		config.Listen = listen
	}

	if port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return err
		}

		config.Port = uint(p)
	}

	if remoteWriteEndpoint != "" {
		config.MimirEndpoint = remoteWriteEndpoint
	}

	return nil
}

func checkConfig(config *Config) error {
	if config.Port > 65535 { //nolint
		return fmt.Errorf("port: %d is not valid", config.Port)
	}

	if !strings.HasPrefix(config.MimirEndpoint, "http") {
		return fmt.Errorf("remote_write_endpoint: should start with http/https")
	}

	return nil
}
