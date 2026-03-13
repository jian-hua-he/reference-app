package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jian-hua-he/reference-app/internal/config"
)

func TestLoad_Defaults(t *testing.T) {
	cfg := config.Load()

	assert.Equal(t, "localhost", cfg.DB.Host)
	assert.Equal(t, 5432, cfg.DB.Port)
	assert.Equal(t, "postgres", cfg.DB.User)
	assert.Equal(t, "postgres", cfg.DB.Password)
	assert.Equal(t, "reference_app", cfg.DB.DBName)
	assert.Equal(t, "disable", cfg.DB.SSLMode)
	assert.Equal(t, 8082, cfg.HTTP.Port)
	assert.Equal(t, 50051, cfg.GRPC.Port)
}

func TestLoad_FromEnv(t *testing.T) {
	testCases := map[string]struct {
		Env    map[string]string
		Assert func(t *testing.T, cfg config.Config)
	}{
		"DB_HOST": {
			Env: map[string]string{"DB_HOST": "dbhost"},
			Assert: func(t *testing.T, cfg config.Config) {
				assert.Equal(t, "dbhost", cfg.DB.Host)
			},
		},
		"DB_PORT": {
			Env: map[string]string{"DB_PORT": "5433"},
			Assert: func(t *testing.T, cfg config.Config) {
				assert.Equal(t, 5433, cfg.DB.Port)
			},
		},
		"DB_USER": {
			Env: map[string]string{"DB_USER": "myuser"},
			Assert: func(t *testing.T, cfg config.Config) {
				assert.Equal(t, "myuser", cfg.DB.User)
			},
		},
		"DB_PASSWORD": {
			Env: map[string]string{"DB_PASSWORD": "secret"},
			Assert: func(t *testing.T, cfg config.Config) {
				assert.Equal(t, "secret", cfg.DB.Password)
			},
		},
		"DB_NAME": {
			Env: map[string]string{"DB_NAME": "mydb"},
			Assert: func(t *testing.T, cfg config.Config) {
				assert.Equal(t, "mydb", cfg.DB.DBName)
			},
		},
		"DB_SSLMODE": {
			Env: map[string]string{"DB_SSLMODE": "require"},
			Assert: func(t *testing.T, cfg config.Config) {
				assert.Equal(t, "require", cfg.DB.SSLMode)
			},
		},
		"HTTP_PORT": {
			Env: map[string]string{"HTTP_PORT": "9090"},
			Assert: func(t *testing.T, cfg config.Config) {
				assert.Equal(t, 9090, cfg.HTTP.Port)
			},
		},
		"GRPC_PORT": {
			Env: map[string]string{"GRPC_PORT": "50052"},
			Assert: func(t *testing.T, cfg config.Config) {
				assert.Equal(t, 50052, cfg.GRPC.Port)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			for k, v := range tc.Env {
				t.Setenv(k, v)
			}
			cfg := config.Load()
			tc.Assert(t, cfg)
		})
	}
}
