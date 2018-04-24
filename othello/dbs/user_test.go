package dbs_test

import (
	"testing"

	"github.com/sdbx/othello-server/othello/dbs"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	dbs.Clear()
	a := assert.New(t)

	// add user
	user := dbs.User{
		Name:   "asdf",
		UserID: "asdf",
	}
	err := dbs.AddUser(&user)
	a.Nil(err)
	sec := user.Secret

	// getuser
	user, err = dbs.GetUserBySecret(sec)
	a.Nil(err)
	a.Equal(user.Name, "asdf", "should equal")

	user, err = dbs.GetUserByUserID("asdf")
	a.Nil(err)
	a.Equal(user.Name, "asdf", "should equal")

	// duplicate userid
	user = dbs.User{
		Name:   "asdf",
		UserID: "asdf",
	}
	err = dbs.AddUser(&user)
	a.NotNil(err)

	// duplicate user name
	user = dbs.User{
		Name:   "asdf",
		UserID: "asdf2",
	}
	err = dbs.AddUser(&user)
	a.Nil(err)

	user, err = dbs.GetUserByUserID("asdf2")
	a.Nil(err)
	a.Equal(user.Name, "asdf0", "should equal")
}
