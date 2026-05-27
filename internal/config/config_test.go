package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setEnv(t *testing.T, kv map[string]string) {
	t.Helper()
	for k, v := range kv {
		t.Setenv(k, v)
	}
}

func validEnv() map[string]string {
	return map[string]string{
		"HTTP_ADDR":              ":12340",
		"DATABASE_URL":           "postgres://app:app@localhost:5432/app?sslmode=disable",
		"KEYCLOAK_URL":           "http://localhost:8080",
		"KEYCLOAK_REALM":         "app",
		"KEYCLOAK_CLIENT_ID":     "app",
		"KEYCLOAK_CLIENT_SECRET": "secret",
	}
}

func TestLoad_Valid(t *testing.T) {
	setEnv(t, validEnv())

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, ":12340", cfg.HTTPAddr)
	assert.Equal(t, "app", cfg.KeycloakRealm)
	assert.Equal(t, "secret", cfg.KeycloakClientSecret)
}

func TestLoad_MissingRequired(t *testing.T) {
	env := validEnv()
	delete(env, "DATABASE_URL")
	setEnv(t, env)

	_, err := Load()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "DATABASE_URL")
}

func TestLoad_InvalidURL(t *testing.T) {
	env := validEnv()
	env["DATABASE_URL"] = "not-a-url"
	setEnv(t, env)

	_, err := Load()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "DATABASE_URL")
}

// envKeys must stay in sync with the mapstructure tags on Config.
func TestEnvKeys(t *testing.T) {
	assert.ElementsMatch(t, []string{
		"HTTP_ADDR",
		"DATABASE_URL",
		"KEYCLOAK_URL",
		"KEYCLOAK_REALM",
		"KEYCLOAK_CLIENT_ID",
		"KEYCLOAK_CLIENT_SECRET",
	}, envKeys())
}
