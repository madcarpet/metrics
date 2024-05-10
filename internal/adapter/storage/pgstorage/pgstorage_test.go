package pgstorage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPGStorage(t *testing.T) {
	//Testing error when connections params are invalid
	db, err := NewPGStorage("test")
	assert.Nil(t, err)
	defer db.Close()
	err = db.IsConnected(context.Background())
	assert.ErrorContains(t, err, "cannot parse")
	//Testing DB timeout
	db2, err := NewPGStorage("dbname=test user=test password=test host=10.89.0.19 port=5432")
	assert.Nil(t, err)
	defer db2.Close()
	err = db2.IsConnected(context.Background())
	assert.ErrorContains(t, err, "failed to connect")
}
