package metrics

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/madcarpet/metrics/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mocks.NewMockRepository(ctrl)

	s.EXPECT().IsConnected(context.Background()).Return(nil)

	pingSvc := NewPingSvc(s)
	err := pingSvc.Ping(context.Background())
	assert.Nil(t, err)

}
