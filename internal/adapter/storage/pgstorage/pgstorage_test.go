package pgstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPGStorage(t *testing.T) {
	//Testing error when connections params are invalid
	_, err := NewPGStorage("test")
	assert.ErrorContains(t, err, "failed to parse as DSN")
	//Testing DB timeout
	_, err = NewPGStorage("dbname=test user=test password=test host=10.89.0.19 port=5432")
	assert.ErrorContains(t, err, "dial error")
}
