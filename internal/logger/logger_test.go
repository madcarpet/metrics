package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInitLogger(t *testing.T) {
	err := Initialize("debug")
	assert.NoError(t, err)

	assert.NotEqual(t, zap.NewNop(), Log)

}
