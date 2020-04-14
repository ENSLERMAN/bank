package store_test

import (
	"github.com/ENSLERMAN/soft-eng/project/internal/app/model"
	"github.com/ENSLERMAN/soft-eng/project/internal/app/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL)

	defer teardown("clients")

	u, err := s.User().Create(model.TestUser(t))
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByLogin(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL)
	defer teardown("clients")

	login := "anakonda3000"
	_, err := s.User().FindByLogin(login)
	assert.Error(t, err)

	u := model.TestUser(t)
	u.Login = login
	s.User().Create(u)

	u, err = s.User().FindByLogin(login)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
