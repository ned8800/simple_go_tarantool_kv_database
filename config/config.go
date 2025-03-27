package config

import (
	"net/http"
	"os"
	"path/filepath"
	"simple_go_tarantool_kv_database/config/defaults"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	ErrInitializeConfig  = "Error initializing config"
	ErrUnmarshalConfig   = "Error unmarshalling config"
	ErrReadConfig        = "Error reading config"
	ErrReadEnvironment   = "Error reading .env file"
	ErrGetDirectory      = "Error getting directory"
	ErrDirectoryNotFound = "Error finding directory"
)

const (
	MaxFindingEnvDepth = 100
)

// server constants
const (
	Address         = "localhost"
	Port            = 8080
	ReadTimeout     = time.Second * 5
	WriteTimeout    = time.Second * 5
	ShutdownTimeout = time.Second * 30
	IdleTimeout     = time.Second * 60
)

// cookie constants
const (
	SessionName   = "session_id"
	SessionLength = 32
	HTTPOnly      = true
	Secure        = false
	SameSite      = http.SameSiteStrictMode
	Path          = "/"
	ExpirationAge = -1
)

type Config struct {
	Server Server `yaml:"server" mapstructure:"server"`
	Cookie Cookie `yaml:"cookie" mapstructure:"cookie"`
}

type Server struct {
	Address         string        `yaml:"address" mapstructure:"address"`
	Port            int           `yaml:"port" mapstructure:"port"`
	ReadTimeout     time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" mapstructure:"shutdown_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" mapstructure:"idle_timeout"`
}

type Cookie struct {
	SessionName   string        `yaml:"session_name" mapstructure:"session_name"`
	SessionLength int           `yaml:"session_length" mapstructure:"session_length"`
	HTTPOnly      bool          `yaml:"http_only" mapstructure:"http_only"`
	Secure        bool          `yaml:"secure" mapstructure:"secure"`
	SameSite      http.SameSite `yaml:"same_site" mapstructure:"same_site"`
	Path          string        `yaml:"path" mapstructure:"path"`
	ExpirationAge int           `yaml:"expiration_age" mapstructure:"expiration_age"`
}

func New() (*Config, error) {
	log.Info().Msg("Initializing config")

	if err := setupViper(); err != nil {
		log.Error().Err(errors.Wrap(err, ErrInitializeConfig)).Msg(errors.Wrap(err, ErrInitializeConfig).Error())
		return nil, errors.Wrap(err, ErrInitializeConfig)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Error().Err(errors.Wrap(err, ErrUnmarshalConfig)).Msg(errors.Wrap(err, ErrUnmarshalConfig).Error())
		return nil, errors.Wrap(err, ErrUnmarshalConfig)
	}

	log.Info().Msg("Config initialized")
	return &config, nil
}

func setupServer() {
	viper.SetDefault("server.address", defaults.Address)
	viper.SetDefault("server.port", defaults.Port)
	viper.SetDefault("server.read_timeout", defaults.ReadTimeout)
	viper.SetDefault("server.write_timeout", defaults.WriteTimeout)
	viper.SetDefault("server.shutdown_timeout", defaults.ShutdownTimeout)
	viper.SetDefault("server.idle_timeout", defaults.IdleTimeout)
}

func findEnvDir() (string, error) {
	log.Info().Msg("Finding environment dir")
	currentDir, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, ErrGetDirectory)
	}

	for i := 0; i < MaxFindingEnvDepth; i++ {
		path := filepath.Join(currentDir, ".env")
		if _, err := os.Stat(path); err == nil {
			log.Info().Msg("Found .env file")
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return "", errors.Wrap(err, ErrDirectoryNotFound)
		}
		currentDir = parentDir
	}

	return "", errors.Wrap(err, ErrDirectoryNotFound)
}

func setupViper() error {
	log.Info().Msg("Initializing viper")

	envDir, err := findEnvDir()
	if err != nil {
		wrapped := errors.Wrap(err, ErrDirectoryNotFound)
		log.Error().Err(wrapped).Msg(wrapped.Error())
		return wrapped
	}

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(envDir)

	if err := viper.ReadInConfig(); err != nil {
		wrapped := errors.Wrap(err, ErrReadEnvironment)
		log.Error().Err(wrapped).Msg(wrapped.Error())
		return wrapped
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(viper.GetString("VIPER_CONFIG_PATH"))

	setupServer()

	if err := viper.MergeInConfig(); err != nil {
		wrapped := errors.Wrap(err, ErrReadConfig)
		log.Error().Err(wrapped).Msg(wrapped.Error())
		return wrapped
	}

	log.Info().Msg("Viper initialized")
	return nil
}
