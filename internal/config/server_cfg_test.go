package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// func TestDefaultConfig(t *testing.T) {
// 	config, err := NewServerConfig()
// 	require.NoError(t, err)

// 	assert.Equal(t, "localhost:8080", config.ServerAddress)
// 	assert.Equal(t, "info", config.LoggingLevel)
// 	assert.Equal(t, int64(300), config.StoreInterval)
// 	assert.Equal(t, false, config.SyncWrite)
// 	assert.Equal(t, "/tmp/metrics-db.json", config.FilePath)
// 	assert.Equal(t, false, config.IsRestore)
// }

func TestEnvConfig(t *testing.T) {

	os.Setenv("ADDRESS", "localhost:9090")
	os.Setenv("LOGLEVEL", "debug")
	os.Setenv("STORE_INTERVAL", "500")
	os.Setenv("FILE_STORAGE_PATH", "metrics-db.json")
	os.Setenv("RESTORE", "false")
	os.Setenv("DATABASE_DSN", "")
	defer os.Clearenv()

	config, err := NewServerConfig()
	require.NoError(t, err)

	assert.Equal(t, "localhost:9090", config.ServerAddress)
	assert.Equal(t, "debug", config.LoggingLevel)
	assert.Equal(t, int64(500), config.StoreInterval)
	assert.Equal(t, false, config.SyncWrite)
	assert.Equal(t, "metrics-db.json", config.FilePath)
	assert.Equal(t, false, config.IsRestore)
	assert.Equal(t, "", config.DBUrl)
}
