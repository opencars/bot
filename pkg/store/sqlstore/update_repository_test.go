package sqlstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/opencars/bot/pkg/model"
	"github.com/opencars/bot/pkg/store/sqlstore"
)

func TestUpdateRepository_Create(t *testing.T) {
	s, teardown := sqlstore.TestDB(t, conf)
	defer teardown("users", "updates")

	user := model.TestUser(t)
	assert.NoError(t, s.User().Create(user))

	update := model.TestUpdate(t)
	assert.NoError(t, s.Update().Create(update))
}

func TestUpdateRepository_FindByID(t *testing.T) {
	s, teardown := sqlstore.TestDB(t, conf)
	defer teardown("users", "updates")

	user := model.TestUser(t)
	assert.NoError(t, s.User().Create(user))

	update := model.TestUpdate(t)
	assert.NoError(t, s.Update().Create(update))

	actual, err := s.Update().FindByID(update.ID)
	assert.NoError(t, err)
	assert.EqualValues(t, update, actual)
}
