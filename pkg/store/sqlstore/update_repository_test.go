package sqlstore_test

import (
	"github.com/opencars/bot/pkg/domain"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/opencars/bot/pkg/store/sqlstore"
)

func TestUpdateRepository_Create(t *testing.T) {
	s, teardown := sqlstore.TestDB(t, conf)
	defer teardown("users", "updates")

	user := domain.TestUser(t)
	assert.NoError(t, s.User().Create(user))

	update := domain.TestUpdate(t)
	assert.NoError(t, s.Update().Create(update))
}

func TestUpdateRepository_FindByID(t *testing.T) {
	s, teardown := sqlstore.TestDB(t, conf)
	defer teardown("users", "updates")

	user := domain.TestUser(t)
	assert.NoError(t, s.User().Create(user))

	update := domain.TestUpdate(t)
	assert.NoError(t, s.Update().Create(update))

	actual, err := s.Update().FindByID(update.ID)
	assert.NoError(t, err)
	assert.EqualValues(t, update, actual)
}
