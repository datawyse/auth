package internal

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// AppConfig contains the application configuration
type AppConfig struct {
	Host                 string `mapstructure:"HOST" envconfig:"HOST" json:"host,omitempty"  default:"localhost"`
	Port                 string `mapstructure:"PORT" envconfig:"PORT" json:"port,omitempty"  default:"8080"`
	GRPCPort             string `mapstructure:"GRPC_PORT" envconfig:"GRPC_PORT" json:"grpcPort,omitempty"  default:"8081"`
	GRPCHost             string `mapstructure:"GRPC_HOST" envconfig:"GRPC_HOST" json:"grpcHost,omitempty"  default:"localhost"`
	RedisDb              string `mapstructure:"REDIS_DB" envconfig:"REDIS_DB" json:"redis_db,omitempty"`
	RedisURI             string `mapstructure:"REDIS_URI" envconfig:"REDIS_URI" json:"redis_uri,omitempty"`
	RedisPassword        string `mapstructure:"REDIS_PASSWORD" envconfig:"REDIS_PASSWORD" json:"redis_password,omitempty"`
	DatabaseName         string `mapstructure:"DATABASE_NAME" envconfig:"DATABASE_NAME" json:"database_name,omitempty" default:"dataflow"`
	DatabaseURI          string `mapstructure:"DATABASE_URI" envconfig:"DATABASE_URI" json:"databaseUri,omitempty" default:"mongodb://localhost:27017"`
	AppMode              string `mapstructure:"APP_MODE" envconfig:"APP_MODE" json:"appMode,omitempty" default:"development"`
	AppName              string `mapstructure:"APP_MODE" envconfig:"APP_MODE" json:"appName" validate:"required"`
	AppVersion           string `mapstructure:"APP_VERSION" envconfig:"APP_VERSION" json:"appVersion" validate:"required"`
	AppUrl               string `mapstructure:"APP_URL" envconfig:"APP_URL" json:"appUrl" validate:"required"`
	ServiceName          string `mapstructure:"SERVICE_NAME" envconfig:"SERVICE_NAME" json:"serviceName,omitempty" default:"dataflow"`
	RequestTimeout       int    `mapstructure:"REQUEST_TIMEOUT" envconfig:"REQUEST_TIMEOUT" json:"requestTimeout,omitempty" default:"10"`
	ServiceTimeout       int    `mapstructure:"SERVICE_TIMEOUT" envconfig:"SERVICE_TIMEOUT" json:"serviceTimeout,omitempty" default:"10"`
	GRPCServiceTimeout   int    `mapstructure:"GRPC_SERVICE_TIMEOUT" envconfig:"GRPC_SERVICE_TIMEOUT" json:"GRPCServiceTimeout,omitempty" default:"10"`
	DatabaseTimeout      int    `mapstructure:"DATABASE_TIMEOUT" envconfig:"DATABASE_TIMEOUT" json:"databaseTimeout,omitempty" default:"10"`
	ConfigFile           string `mapstructure:"CONFIG_FILE" envconfig:"CONFIG_FILE" json:"configFile,omitempty" default:"config.yml"`
	KeycloakServer       string `mapstructure:"KEYCLOAK_SERVER" json:"keycloakServer,omitempty" envconfig:"KEYCLOAK_SERVER" default:"http://localhost:8080"`
	KeycloakRealm        string `mapstructure:"KEYCLOAK_REALM" json:"keycloakRealm,omitempty" envconfig:"KEYCLOAK_REALM" default:"datawyse"`
	KeycloakClientId     string `mapstructure:"KEYCLOAK_CLIENT_ID" json:"keycloakClientId,omitempty" envconfig:"KEYCLOAK_CLIENT_ID" default:"datawyse"`
	KeycloakClientSecret string `mapstructure:"KEYCLOAK_CLIENT_SECRET" json:"keycloakClientSecret,omitempty" envconfig:"KEYCLOAK_CLIENT_SECRET"`
	CollectorURL         string `mapstructure:"COLLECTOR_URL" json:"collectorUrl,omitempty" envconfig:"COLLECTOR_URL" default:"localhost:4317"`
}

func (appConfig *AppConfig) String() string {
	return fmt.Sprintf(
		"host: %s, port: %s, databaseName: %s, databaseUri: %s, configFile: %s, appMode: %s",
		appConfig.Host,
		appConfig.Port,
		appConfig.DatabaseName,
		appConfig.DatabaseURI,
		appConfig.ConfigFile,
		appConfig.AppMode,
	)
}

func NewConfig(configFile string) (*AppConfig, error) {
	if configFile != "" {
		// read the config file from root directory
		viper.SetConfigType("yml")
		viper.SetConfigFile(configFile)
	}

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AutomaticEnv()

	if configFile != "" {
		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	var appConfig AppConfig
	if err := viper.Unmarshal(&appConfig); err != nil {
		return nil, err
	}

	return &appConfig, nil
}

var AppConfigModule = NewConfig
