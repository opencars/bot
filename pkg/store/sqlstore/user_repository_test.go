package sqlstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/opencars/bot/pkg/model"
	"github.com/opencars/bot/pkg/store/sqlstore"
)

func TestUserRepository_Create(t *testing.T) {
	s, teardown := sqlstore.TestDB(t, conf)
	defer teardown("users")

	user := model.TestUser(t)
	assert.NoError(t, s.User().Create(user))
}

func TestUserRepository_FindByID(t *testing.T) {
	s, teardown := sqlstore.TestDB(t, conf)
	defer teardown("users")

	user := model.TestUser(t)
	assert.NoError(t, s.User().Create(user))

	actual, err := s.User().FindByID(user.ID)
	assert.NoError(t, err)
	assert.EqualValues(t, user, actual)
}
