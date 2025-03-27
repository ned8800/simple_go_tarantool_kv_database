package config

import (
	"net/http"
	"os"
	"path/filepath"
	"simple_go_tarantool_kv_database/config/defaults"
	errorconstants "simple_go_tarantool_kv_database/error_constants"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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
}

type Server struct {
	Address         string        `yaml:"address" mapstructure:"address"`
	Port            int           `yaml:"port" mapstructure:"port"`
	ReadTimeout     time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" mapstructure:"shutdown_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" mapstructure:"idle_timeout"`
}

func New() (*Config, error) {
	log.Info().Msg("Initializing config")

	if err := setupViper(); err != nil {
		log.Error().Err(errors.Wrap(err, errorconstants.ErrInitializeConfig)).Msg(errors.Wrap(err, errorconstants.ErrInitializeConfig).Error())
		return nil, errors.Wrap(err, errorconstants.ErrInitializeConfig)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Error().Err(errors.Wrap(err, errorconstants.ErrUnmarshalConfig)).Msg(errors.Wrap(err, errorconstants.ErrUnmarshalConfig).Error())
		return nil, errors.Wrap(err, errorconstants.ErrUnmarshalConfig)
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
		return "", errors.Wrap(err, errorconstants.ErrGetDirectory)
	}

	for i := 0; i < MaxFindingEnvDepth; i++ {
		path := filepath.Join(currentDir, ".env")
		if _, err := os.Stat(path); err == nil {
			log.Info().Msg("Found .env file")
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return "", errors.Wrap(err, errorconstants.ErrDirectoryNotFound)
		}
		currentDir = parentDir
	}

	return "", errors.Wrap(err, errorconstants.ErrDirectoryNotFound)
}

func setupViper() error {
	log.Info().Msg("Initializing viper")

	envDir, err := findEnvDir()
	if err != nil {
		wrapped := errors.Wrap(err, errorconstants.ErrDirectoryNotFound)
		log.Error().Err(wrapped).Msg(wrapped.Error())
		return wrapped
	}

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(envDir)

	if err := viper.ReadInConfig(); err != nil {
		wrapped := errors.Wrap(err, errorconstants.ErrReadEnvironment)
		log.Error().Err(wrapped).Msg(wrapped.Error())
		return wrapped
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(viper.GetString("VIPER_CONFIG_PATH"))

	setupServer()

	if err := viper.MergeInConfig(); err != nil {
		wrapped := errors.Wrap(err, errorconstants.ErrReadConfig)
		log.Error().Err(wrapped).Msg(wrapped.Error())
		return wrapped
	}

	log.Info().Msg("Viper initialized")
	return nil
}
