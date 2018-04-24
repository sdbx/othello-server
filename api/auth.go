package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/sdbx/othello-server/othello/dbs"
)

func handleAuth(w http.ResponseWriter, r *http.Request, user goth.User) {
	user2, err := dbs.GetUserByUserID(user.UserID)
	if err != nil {
		user2 = dbs.User{
			Name:    user.NickName,
			UserID:  user.UserID,
			Profile: user.AvatarURL,
		}
		err = dbs.AddUser(&user2)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
	}
	http.Redirect(w, r, os.Getenv("AUTH_CALLBACK")+"?secret="+user2.Secret, 301)
}
func authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	handleAuth(w, r, user)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if user, err := gothic.CompleteUserAuth(w, r); err == nil {
		handleAuth(w, r, user)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}
