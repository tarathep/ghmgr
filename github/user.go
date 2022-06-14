package github

import (
	"encoding/json"
	"errors"
	"strings"

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

func (user User) EmailToUsername(caches []model.Cache, email string) string {
	for _, c := range caches {
		if c.Email == email {
			return c.Username
		}
	}
	return ""
}
func (user User) UsernameToEmail(caches []model.Cache, username string) string {
	for _, c := range caches {
		if c.Username == username {
			return c.Email
		}
	}
	return ""
}
func (user User) CheckAlreadyMemberByEmail(caches []model.Cache, email string) bool {
	for _, c := range caches {
		if c.Email == email {
			return true
		}
	}
	return false
}

func (user User) CheckAlreadyMemberTeamByEmail(caches []model.Cache, email string, team string) bool {
	for _, c := range caches {
		if c.Email == email && isCacheTeamAlready(c.Team, team) {
			return true
		}
	}
	return false
}
func isCacheTeamAlready(team string, teams_string string) bool {
	strings.Split(teams_string, "|")
	return false
}

func (user User) CheckEmailInList(emails []string, email string) bool {
	for _, e := range emails {
		if e == email {
			return true
		}
	}
	return false
}
