package github

import (
	"encoding/json"
	"errors"

	"github.com/tarathep/ghmgr/login"
	"github.com/tarathep/ghmgr/model"
)

type User struct {
	Auth login.Auth
}

func (user User) UserInfo(username string) (error, model.User) {
	github := GitHub{Auth: user.Auth}
	_, statusCode, bodyBytes := github.Request("GET", "https://api.github.com/users/"+username, nil)

	if statusCode != 200 {
		return errors.New(github.GetMessage(bodyBytes)), model.User{}
	}

	usr := model.User{}
	json.Unmarshal(bodyBytes, &usr)

	return nil, usr
}
