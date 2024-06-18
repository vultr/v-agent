// Package config ensures the app is configured properly
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/vultr/v-agent/pkg/connectors"
	"github.com/vultr/v-agent/pkg/util"
	"github.com/vultr/v-agent/pkg/wrkld"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/util/validation"
)

var cfg Config
var labels map[string]string

// Config is the CLI options wrapped in a struct
type Config struct {
	Name       string
	ConfigFile string

	Version string

	Debug         bool              `yaml:"debug"`
	Interval      uint              `yaml:"interval"`
	Endpoint      string            `yaml:"endpoint"`
	BasicAuthUser string            `yaml:"basic_auth_user"`
	BasicAuthPass string            `yaml:"basic_auth_pass"`
	LabelsConfig  map[string]string `yaml:"labels_config"`
	ProbesAPI     ProbesAPI         `yaml:"probes_api"`
	MetricsConfig MetricsConfig     `yaml:"metrics_config"`

	zapConfig *zap.Config
	zapLogger *zap.Logger
	zapSugar  *zap.SugaredLogger
}

// MetricsConfig contains metrics configuration
type MetricsConfig struct {
	Agent      AgentMetrics      `yaml:"agent"`
	Kubernetes KubernetesMetrics `yaml:"kubernetes"`
}

// AgentMetrics metrics that are collected when ran as an agent (daemon)
type AgentMetrics struct {
	LoadAvg      LoadAvg      `yaml:"load_avg"`
	CPU          CPU          `yaml:"cpu"`
	Memory       Memory       `yaml:"memory"`
	NIC          NIC          `yaml:"nic"`
	DiskStats    DiskStats    `yaml:"disk_stats"`
	Filesystem   Filesystem   `yaml:"file_system"`
	Kubernetes   Kubernetes   `yaml:"kubernetes"`
	Konnectivity Konnectivity `yaml:"konnectivity"`
	Etcd         Etcd         `yaml:"etcd"`
	NginxVTS     NginxVTS     `yaml:"nginx_vts"`
	VCDNAgent    VCDNAgent    `yaml:"v_cdn_agent"`
	HAProxy      HAProxy      `yaml:"haproxy"`
	Ceph         Ceph         `yaml:"ceph"`
	VDNS         VDNS         `yaml:"v_dns"`
	SMART        SMART        `yaml:"smart"`
}

// KubernetesMetrics metrics that are collected when ran as an operator (in k8s)
type KubernetesMetrics struct {
	Pods Pods `yaml:"pods"`
	DCGM DCGM `yaml:"dcgm"`
}

// ProbesAPI probes API definition
type ProbesAPI struct {
	Listen string `yaml:"listen"`
	Port   uint   `yaml:"port"`
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

// SMART config
type SMART struct {
	Enabled      bool     `yaml:"enabled"`
	BlockDevices []string `yaml:"block_devices"`
}

// Pods config
type Pods struct {
	Enabled    bool     `yaml:"enabled"`
	Namespaces []string `yaml:"namespaces"`
}

// DCGM config
type DCGM struct {
	Enabled   bool   `yaml:"enabled"`
	Namespace string `yaml:"namespace"`
	Endpoint  string `yaml:"endpoint"`
}

// NewConfig returns a Config struct that can be used to reference configuration
// NewConfig does the following:
//   - Runs initCLI (sets and read CLI switches)
//   - Runs initConfig (reads config from files)
func NewConfig(name, version string) (*Config, error) {
	// Stage 1: Setup CLI flags
	initCLI(&cfg)

	// Stage 2: Setup config file
	if err := initConfig(&cfg); err != nil {
		return nil, err
	}

	// Stage 3: initialize logging
	initLogging(&cfg)

	// Stage 4: Setup env vars
	if err := initEnv(&cfg); err != nil {
		return nil, err
	}

	// Stage 5: Initialize labels
	if err := initLabels(&cfg); err != nil {
		return nil, err
	}

	// Stage 6: Check configuration
	if err := checkConfig(&cfg); err != nil {
		return nil, err
	}

	cfg.Name = name
	cfg.Version = version

	return &cfg, nil
}

// initCLI initializes CLI switches
func initCLI(config *Config) {
	flag.BoolVar(&config.Debug, "debug", true, "Debug output")
	flag.StringVar(&config.ProbesAPI.Listen, "listen", "127.0.0.1", "Listen address")
	flag.UintVar(&config.ProbesAPI.Port, "port", 7091, "Listen port") //nolint
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
		config.ProbesAPI.Listen = listen
	}

	if port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return err
		}

		config.ProbesAPI.Port = uint(p)
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
		config.MetricsConfig.Agent.Kubernetes.Kubeconfig = kubeconfig
	}

	return nil
}

func checkConfig(config *Config) error {
	log := zap.L().Sugar()

	if config.ProbesAPI.Port > 65535 { //nolint
		return fmt.Errorf("%d: %w", config.ProbesAPI.Port, ErrPortInvalid)
	}

	if config.Interval < 1 {
		return fmt.Errorf("%d: %w", config.ProbesAPI.Port, ErrIntervalInvalid)
	}

	if !strings.HasPrefix(config.Endpoint, "http") {
		return fmt.Errorf("remote_write_endpoint: %w", ErrMissingScheme)
	}

	if config.MetricsConfig.Kubernetes.Pods.Enabled {
		// try to get k8s connection, if error return
		if !inK8s() {
			return ErrNotInK8s
		}

		for i := range config.MetricsConfig.Kubernetes.Pods.Namespaces {
			errs := validation.IsValidLabelValue(config.MetricsConfig.Kubernetes.Pods.Namespaces[i])
			if len(errs) > 0 {
				log.Error(errs)

				return fmt.Errorf("pods: %w", ErrKubernetesNamespaceInvalid)
			}
		}
	}

	if config.MetricsConfig.Agent.SMART.Enabled {
		if len(config.MetricsConfig.Agent.SMART.BlockDevices) > 0 {
			for i := range config.MetricsConfig.Agent.SMART.BlockDevices {
				if _, err := os.Stat(config.MetricsConfig.Agent.SMART.BlockDevices[i]); errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("%s: %w", config.MetricsConfig.Agent.SMART.BlockDevices[i], ErrSMARTDeviceNotExist)
				}
			}
		}
	}

	if config.MetricsConfig.Kubernetes.DCGM.Enabled {
		if !inK8s() {
			return ErrNotInK8s
		}

		if config.MetricsConfig.Kubernetes.DCGM.Namespace == "" {
			return fmt.Errorf("dcgm.namespace: %w", ErrKubernetesNamespaceNotSet)
		}

		clientset, err := connectors.GetKubernetesConn()
		if err != nil {
			return err
		}

		if !wrkld.NamespaceExists(clientset, config.MetricsConfig.Kubernetes.DCGM.Namespace) {
			return fmt.Errorf("dcgm.namespace: %w", ErrKubernetesNamespaceNotExist)
		}

		if config.MetricsConfig.Kubernetes.DCGM.Endpoint == "" {
			return ErrDCGMEndpointNotSet
		}

		if !wrkld.EndpointExists(clientset, config.MetricsConfig.Kubernetes.DCGM.Namespace, config.MetricsConfig.Kubernetes.DCGM.Endpoint) {
			return ErrDCGMEndpointNotExist
		}
	}

	return nil
}

func inK8s() bool {
	v := os.Getenv("KUBERNETES_SERVICE_HOST")

	return v != ""
}
