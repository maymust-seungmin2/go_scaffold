// Package config provides configuration management for the application.
package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Config struct {
	HTTPAddr    string `mapstructure:"HTTP_ADDR"    validate:"required"`
	DatabaseURL string `mapstructure:"DATABASE_URL" validate:"required,url"`

	KeycloakURL          string `mapstructure:"KEYCLOAK_URL"           validate:"required,url"`
	KeycloakRealm        string `mapstructure:"KEYCLOAK_REALM"         validate:"required"`
	KeycloakClientID     string `mapstructure:"KEYCLOAK_CLIENT_ID"     validate:"required"`
	KeycloakClientSecret string `mapstructure:"KEYCLOAK_CLIENT_SECRET" validate:"required"`
}

// Load reads configuration from environment variables and validates it.
func Load() (*Config, error) {
	v := viper.New()

	for _, key := range envKeys() {
		if err := v.BindEnv(key); err != nil {
			return nil, oops.In("config").With("env_key", key).Wrapf(err, "bind env var")
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, oops.In("config").Wrapf(err, "unmarshal config")
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, oops.
			In("config").
			Code("invalid_config").
			Hint("Set the required environment variables before starting the server.").
			Wrapf(err, "invalid configuration:\n%s", formatValidationError(err))
	}

	return &cfg, nil
}

// validateConfig reports field violations using the env var names (mapstructure
// tags) rather than Go field names, so the message points at what operators set.
func validateConfig(cfg *Config) error {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(f reflect.StructField) string {
		return f.Tag.Get("mapstructure")
	})
	return validate.Struct(cfg)
}

// envKeys returns the environment variable names bound to Config, read from the
// mapstructure tags so the struct remains the single source of truth.
func envKeys() []string {
	t := reflect.TypeFor[Config]()
	keys := make([]string, 0, t.NumField())
	for f := range t.Fields() {
		if tag := f.Tag.Get("mapstructure"); tag != "" {
			keys = append(keys, tag)
		}
	}
	return keys
}

func formatValidationError(err error) string {
	var verr validator.ValidationErrors
	if !errors.As(err, &verr) {
		return err.Error()
	}

	lines := make([]string, 0, len(verr))
	for _, fe := range verr {
		lines = append(lines, fmt.Sprintf("  - %s: failed rule %q", fe.Field(), fe.Tag()))
	}
	return strings.Join(lines, "\n")
}
