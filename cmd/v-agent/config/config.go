// Package config ensures the app is configured properly
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/vultr/v-agent/cmd/v-agent/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

var (
	// ErrNotVultrVendor returned if the bios manufacturer is not "Vultr"
	ErrNotVultrVendor = errors.New("not vultr vendor")
)

var config *Config

// Config is the CLI options wrapped in a struct
type Config struct {
	ConfigFile string

	Version string

	Debug         bool          `yaml:"debug"`
	Listen        string        `yaml:"listen"`
	Port          uint          `yaml:"port"`
	Interval      uint          `yaml:"interval"`
	SubID         string        `yaml:"subid"`
	Product       string        `yaml:"product"`
	Endpoint      string        `yaml:"endpoint"`
	BasicAuthUser string        `yaml:"basic_auth_user"`
	BasicAuthPass string        `yaml:"basic_auth_pass"`
	MetricsConfig MetricsConfig `yaml:"metrics_config"`

	zapConfig *zap.Config
	zapLogger *zap.Logger
	zapSugar  *zap.SugaredLogger
}

// MetricsConfig contains metrics configuration
type MetricsConfig struct {
	LoadAvg    LoadAvg    `yaml:"load_avg"`
	CPU        CPU        `yaml:"cpu"`
	Memory     Memory     `yaml:"memory"`
	NIC        NIC        `yaml:"nic"`
	DiskStats  DiskStats  `yaml:"disk_stats"`
	Filesystem Filesystem `yaml:"file_system"`
	Kubernetes Kubernetes `yaml:"kubernetes"`
	Etcd       Etcd       `yaml:"etcd"`
}

// LoadAvg configuration
type LoadAvg struct {
	Enabled bool `yaml:"enabled"`
}

// CPU configuration
type CPU struct {
	Enabled bool `yaml:"enabled"`
}

// Memory configuration
type Memory struct {
	Enabled bool `yaml:"enabled"`
}

// NIC configuration
type NIC struct {
	Enabled bool `yaml:"enabled"`
}

// DiskStats configuration
type DiskStats struct {
	Enabled bool   `yaml:"enabled"`
	Filter  string `yaml:"filter"`
}

// Filesystem configuration
type Filesystem struct {
	Enabled bool `yaml:"enabled"`
}

// Kubernetes config
type Kubernetes struct {
	Enabled    bool   `yaml:"enabled"`
	Endpoint   string `yaml:"endpoint"`
	Kubeconfig string `yaml:"kubeconfig"`
}

// Etcd config
type Etcd struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
	CACert   string `yaml:"cacert"`
	Cert     string `yaml:"cert"`
	Key      string `yaml:"key"`
}

// NewConfig returns a Config struct that can be used to reference configuration
// NewConfig does the following:
//   - Runs initCLI (sets and read CLI switches)
//   - Runs initConfig (reads config from files)
func NewConfig(name, version string) (*Config, error) {
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

	config.Version = version

	return config, nil
}

// initCLI initializes CLI switches
func initCLI(config *Config) {
	flag.BoolVar(&config.Debug, "debug", true, "Debug output")
	flag.StringVar(&config.Listen, "listen", "127.0.0.1", "Listen address")
	flag.UintVar(&config.Port, "port", 7091, "Listen port") //nolint
	flag.StringVar(&config.ConfigFile, "config", "./config.yaml", "Path for the config.yaml configuration file")
	flag.UintVar(&config.Interval, "interval", 60, "Metrics gather interval") //nolint
	flag.StringVar(&config.SubID, "subid", "", "Subid")
	flag.StringVar(&config.Product, "product", "", "Product")
	flag.StringVar(&config.Endpoint, "endpoint", "http://localhost:8080", "Endpoint to remotely write metrics to")
	flag.StringVar(&config.BasicAuthUser, "basic-auth-user", "", "Basic auth user")
	flag.StringVar(&config.BasicAuthPass, "basic-auth-pass", "", "Basic auth password")
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
	interval := os.Getenv("INTERVAL")
	subid := os.Getenv("SUBID")
	product := os.Getenv("PRODUCT")
	endpoint := os.Getenv("ENDPOINT")
	basicAuthUser := os.Getenv("BASIC_AUTH_USER")
	basicAuthPass := os.Getenv("BASIC_AUTH_PASS")
	kubeconfig := os.Getenv("KUBECONFIG")

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

	if interval != "" {
		p, err := strconv.Atoi(interval)
		if err != nil {
			return err
		}

		config.Interval = uint(p)
	}

	if subid != "" {
		config.SubID = subid
	}

	if product != "" {
		config.Product = product
	}

	if endpoint != "" {
		config.Endpoint = endpoint
	}

	if basicAuthUser != "" {
		config.BasicAuthUser = basicAuthUser
	}

	if basicAuthPass != "" {
		config.BasicAuthPass = basicAuthPass
	}

	if kubeconfig != "" {
		config.MetricsConfig.Kubernetes.Kubeconfig = kubeconfig
	}

	return nil
}

func checkConfig(config *Config) error {
	log := zap.L().Sugar()

	vendor, err := util.GetBIOSVendor()
	if err != nil {
		return err
	}

	if *vendor != "Vultr" {
		return ErrNotVultrVendor
	}

	if config.Port > 65535 { //nolint
		return fmt.Errorf("port: %d is not valid", config.Port)
	}

	if config.Interval < 1 {
		return fmt.Errorf("interval: %d is not valid", config.Interval)
	}

	if config.SubID == "" {
		log.Info("subid is not set, pulling from metadata server")

		subid, err := util.GetSubID()
		if err != nil {
			return err
		}

		log.Infof("setting subid to %s", *subid)
		config.SubID = *subid
	}

	if config.Product == "" {
		log.Info("product is not set, setting to \"unknown\"")

		config.Product = "unknown"
	}

	if !strings.HasPrefix(config.Endpoint, "http") {
		return fmt.Errorf("remote_write_endpoint: should start with http/https")
	}

	return nil
}
