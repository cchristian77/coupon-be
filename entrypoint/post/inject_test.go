package post

import (
	m "base_project/mock"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewController_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := m.NewMockRepository(ctrl)
	writerDB, _, err := m.NewMockDB()
	if err != nil {
		t.Fatal(err)
	}

	postController, err := NewController(context.Background(), repo, writerDB)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, postController)
}
