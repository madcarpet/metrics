package pgstorage

import (
	"context"
	"fmt"
	"testing"

	"github.com/madcarpet/metrics/internal/entity"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPGStorageCont(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "metricsdb",
		},
		// WaitingFor: wait.ForLog("database system is ready to accept connections"),
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer postgresContainer.Terminate(ctx)

	host, err := postgresContainer.Host(ctx)
	require.NoError(t, err)
	port, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := fmt.Sprintf("postgres://test:test@%s:%s/metricsdb?sslmode=disable", host, port.Port())
	t.Log("Connecting to:", dsn)
	db, err := NewPGStorage(dsn)
	require.NoError(t, err)
	defer db.Close()
	// Test IsConnected.
	err = db.IsConnected(ctx)
	require.NoError(t, err)
	// Test Migration.
	err = DBMigration("../../../../migrations", db.DB)
	require.NoError(t, err)
	// Test Update.
	gMet := entity.Metric{
		Type:  entity.Gauge,
		Name:  "gaugeMetric",
		Value: 32.254,
	}
	cMet := entity.Metric{
		Type:  entity.Counter,
		Name:  "counterMetric",
		Value: 100,
	}
	err = db.UpdateMetric(ctx, gMet)
	require.NoError(t, err)
	err = db.UpdateMetric(ctx, cMet)
	require.NoError(t, err)
	// Test Updates.
	metricSet := []entity.Metric{gMet, cMet}
	err = db.UpdateMetrics(ctx, metricSet)
	require.NoError(t, err)
	//Test Get by Name and Type.
	gotMetric, err := db.GetByNameAndType(ctx, "counterMetric", entity.Counter)
	require.NoError(t, err)
	require.Equal(t, float64(200), gotMetric.Value)
	//Test Get all metrics.
	gotMetrics, err := db.GetAllMetrics(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(gotMetrics))
}
