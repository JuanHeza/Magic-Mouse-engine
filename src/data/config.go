package data

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

// NewConfig returnsCoreWsagnew Config
func NewConfig() *Config {
	return &Config{}
}

// Config application configuration
type Config struct {
	//Logging Log Params
	Logging struct {
		Level string `yaml:"level" envconfig:"LOGGING_LEVEL"`
		Sink  string `yaml:"sink" envconfig:"LOGGING_SINK"`
	} `yaml:"logging"`

	//Server Server Params
	Server struct {
		//ListenAddress Protocol and Address to listen on http://0.0.0.0:18083, https://0.0.0.0:18083 etc
		ListenAddress string `yaml:"listenAddress" envconfig:"optional,SERVER_LISTEN_ADDRESS"`
		//IngressCIDR   []string `yaml:"ingressCIDR" envconfig:"SERVER_INGRESS_CIDR"`
	} `yaml:"server"`
	//DB Database connection
	DB struct {
		Server   string `yaml:"server" envconfig:"optional,DB_SERVER"`
		Port     int    `yaml:"port" envconfig:"optional,DB_PORT"`
		User     string `yaml:"user" envconfig:"optional,DB_USER"`
		Password string `yaml:"password" envconfig:"optional,DB_PASSWORD"`
		//ConnectionString ODBC Connection String
		ConnectionString     string        `yaml:"connectionString" envconfig:"optional,DB_CONNECTION_STRING"`
		InitConnectRetries   int           `yaml:"initConnectRetries" envconfig:"optional,DB_INIT_CONNECT_RETRIES"`
		InitConnectLimitLog  int           `yaml:"initConnectLimitLog" envconfig:"optional,DB_INIT_CONNECT_LIMIT_LOG"`
		InitConnectSleep     time.Duration `yaml:"initConnectSleep" envconfig:"optional,DB_INIT_CONNECT_SLEEP"`
		SQLReconnectOnErrors []string      `yaml:"sqlReconnectOnErrors" envconfig:"optional,DB_SQL_RECONNECT_ON_ERRORS"`
	} `yaml:"db"`
}

// Init initialize
func (cfg *Config) Init(filePath string) error {
	println("File:", filePath)
	err := processFile(filePath, cfg)
	if err != nil {
		println("Error processing file:", err.Error())
		return err
	}
	return nil
}

func processFile(filePath string, cfg *Config) error {
	newVar := filepath.Clean(filePath)
	f, err := os.Open(newVar)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			println("Error closing file:", err.Error())
		}
	}()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return err
	}
	return nil
}

