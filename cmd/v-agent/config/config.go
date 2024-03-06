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
	"github.com/vultr/v-agent/spec/connectors"
	"k8s.io/apimachinery/pkg/util/validation"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

var (
	// ErrNotVultrVendor returned if the bios manufacturer is not "Vultr"
	ErrNotVultrVendor = errors.New("not vultr vendor")
)

var config *Config
var labels map[string]string

// Config is the CLI options wrapped in a struct
type Config struct {
	Name       string
	ConfigFile string

	Version string

	Debug         bool              `yaml:"debug"`
	Listen        string            `yaml:"listen"`
	Port          uint              `yaml:"port"`
	Interval      uint              `yaml:"interval"`
	Endpoint      string            `yaml:"endpoint"`
	BasicAuthUser string            `yaml:"basic_auth_user"`
	BasicAuthPass string            `yaml:"basic_auth_pass"`
	CheckVendor   bool              `yaml:"check_vendor"`
	LabelsConfig  map[string]string `yaml:"labels_config"`
	MetricsConfig MetricsConfig     `yaml:"metrics_config"`

	zapConfig *zap.Config
	zapLogger *zap.Logger
	zapSugar  *zap.SugaredLogger
}

// MetricsConfig contains metrics configuration
type MetricsConfig struct {
	LoadAvg        LoadAvg        `yaml:"load_avg"`
	CPU            CPU            `yaml:"cpu"`
	Memory         Memory         `yaml:"memory"`
	NIC            NIC            `yaml:"nic"`
	DiskStats      DiskStats      `yaml:"disk_stats"`
	Filesystem     Filesystem     `yaml:"file_system"`
	Kubernetes     Kubernetes     `yaml:"kubernetes"`
	Konnectivity   Konnectivity   `yaml:"konnectivity"`
	Etcd           Etcd           `yaml:"etcd"`
	NginxVTS       NginxVTS       `yaml:"nginx_vts"`
	VCDNAgent      VCDNAgent      `yaml:"v_cdn_agent"`
	HAProxy        HAProxy        `yaml:"haproxy"`
	Ganesha        Ganesha        `yaml:"ganesha"`
	Ceph           Ceph           `yaml:"ceph"`
	VDNS           VDNS           `yaml:"v_dns"`
	KubernetesPods KubernetesPods `yaml:"kubernetes_pods"`
	SMART          SMART          `yaml:"smart"`
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

// Konnectivity config
type Konnectivity struct {
	Enabled         bool   `yaml:"enabled"`
	MetricsEndpoint string `yaml:"metrics_endpoint"`
	HealthEndpoint  string `yaml:"health_endpoint"`
}

// Etcd config
type Etcd struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
	CACert   string `yaml:"cacert"`
	Cert     string `yaml:"cert"`
	Key      string `yaml:"key"`
}

// NginxVTS config
type NginxVTS struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
}

// VCDNAgent config
type VCDNAgent struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
}

// HAProxy config
type HAProxy struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
}

// Ganesha config
type Ganesha struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
}

// Ceph config
type Ceph struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
}

// VDNS config
type VDNS struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
}

// KubernetesPods config
type KubernetesPods struct {
	Enabled    bool     `yaml:"enabled"`
	Namespaces []string `yaml:"namespaces"`
}

// SMART config
type SMART struct {
	Enabled      bool     `yaml:"enabled"`
	BlockDevices []string `yaml:"block_devices"`
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

		// Stage 5: Initialize labels
		if err := initLabels(config); err != nil {
			return nil, err
		}

		// Stage 6: Check configuration
		if err := checkConfig(config); err != nil {
			return nil, err
		}
	}

	config.Name = name
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

// initLabels initializes labels to minimize metadata traffic (if any)
func initLabels(config *Config) error {
	log := zap.L().Sugar()

	product := make(map[string]string)
	for k, v := range config.LabelsConfig {
		if k == "product" {
			product[k] = v
		}
	}

	if labels == nil {
		labels = make(map[string]string)
	}

	for k, v := range config.LabelsConfig {
		switch k {
		case "hostname":
			if v == "" {
				labels[k] = os.Getenv("HOSTNAME")
				if labels[k] == "" {
					hostname, err := os.Hostname()
					if err != nil {
						return err
					}

					labels[k] = hostname
				}
			} else {
				labels[k] = v
			}

		case "subid":
			if v == "" {
				if product["product"] == "" {
					log.Warn("product label unset, unable to determine subid")
					continue
				}

				subid, err := util.GetSubID(product["product"])
				if err != nil {
					return err
				}

				labels[k] = *subid
			} else {
				labels[k] = v
			}
		case "vpsid":
			if v == "" {
				vpsid, err := util.GetVPSID()
				if err != nil {
					return err
				}

				labels[k] = *vpsid
			} else {
				labels[k] = v
			}
		default:
			if v == "" {
				log.Warnf("label %q is empty, will be omitted", k)
				continue
			}

			labels[k] = v
		}
	}

	return nil
}

func initEnv(config *Config) error {
	listen := os.Getenv("LISTEN")
	port := os.Getenv("PORT")
	interval := os.Getenv("INTERVAL")
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

	if config.CheckVendor {
		vendor, err := util.GetBIOSVendor()
		if err != nil {
			return err
		}

		if *vendor != "Vultr" {
			return ErrNotVultrVendor
		}
	}

	if config.Port > 65535 { //nolint
		return fmt.Errorf("port: %d is not valid", config.Port)
	}

	if config.Interval < 1 {
		return fmt.Errorf("interval: %d is not valid", config.Interval)
	}

	if !strings.HasPrefix(config.Endpoint, "http") {
		return fmt.Errorf("remote_write_endpoint: should start with http/https")
	}

	if config.MetricsConfig.KubernetesPods.Enabled {
		// try to get k8s connection, if error return
		_, err := connectors.GetKubernetesConn()
		if err != nil {
			return err
		}

		for i := range config.MetricsConfig.KubernetesPods.Namespaces {
			errs := validation.IsValidLabelValue(config.MetricsConfig.KubernetesPods.Namespaces[i])
			if len(errs) > 0 {
				log.Error(errs)

				return fmt.Errorf("kubernetes_pods.namespaces invalid")
			}
		}
	}

	if config.MetricsConfig.SMART.Enabled {
		if len(config.MetricsConfig.SMART.BlockDevices) > 0 {
			for i := range config.MetricsConfig.SMART.BlockDevices {
				if _, err := os.Stat(config.MetricsConfig.SMART.BlockDevices[i]); errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("smart.block_devices %s does not exist", config.MetricsConfig.SMART.BlockDevices[i])
				}
			}
		}
	}

	return nil
}
