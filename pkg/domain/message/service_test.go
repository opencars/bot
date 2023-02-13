package message_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/opencars/bot/pkg/domain/message"
	"github.com/opencars/bot/pkg/domain/model"
	"github.com/opencars/bot/pkg/store/mockstore"
)

func TestService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	msg := model.TestMessage(t)

	repo := mockstore.NewMockMessageRepository(ctrl)
	repo.EXPECT().Create(gomock.Any(), msg)

	svc, err := message.NewService(repo)
	require.NoError(t, err)

	assert.NoError(t, svc.Create(context.Background(), msg))
}
