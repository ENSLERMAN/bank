package sqlstore_test

import (
	"github.com/ENSLERMAN/soft-eng/project/internal/app/model"
	"github.com/ENSLERMAN/soft-eng/project/internal/app/store"
	"github.com/ENSLERMAN/soft-eng/project/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

 /*
 ТЕСТЫ МЕТОДОВ ДЛЯ ЮЗЕРА.
  */

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("clients")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByLogin(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	s := sqlstore.New(db)
	defer teardown("clients")

	login := "anakonda3000"
	_, err := s.User().FindByLogin(login)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	u.Login = login
	s.User().Create(u)

	u, err = s.User().FindByLogin(login)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByID(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("clients")

	s := sqlstore.New(db)
	u1 := model.TestUser(t)
	s.User().Create(u1)

	u2, err := s.User().FindByID(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
