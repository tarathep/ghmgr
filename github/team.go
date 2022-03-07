package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/tarathep/githuby/login"
	"github.com/tarathep/githuby/model"
)

type Team struct {
	Auth  login.Auth
	Owner string
}

func (team Team) GetRepos(teamName string, page string) model.Repos {
	github := GitHub{Auth: team.Auth}
	statusCode, bodyBytes := github.Request("GET", "https://api.github.com/orgs/"+team.Owner+"/teams/"+teamName+"/repos?page="+page, nil)

	if statusCode != 200 {
		log.Println(statusCode, github.GetMessage(bodyBytes))
	}

	repos := model.Repos{}

	json.Unmarshal(bodyBytes, &repos)

	return repos
}

func (team Team) GetRepoList(teamName string) []string {
	var nameRepos []string
	total := 0

	for i := 0; true; i++ {
		pagex := strconv.Itoa((i + 1))

		repos := team.GetRepos(teamName, pagex)

		if len(repos) == 0 {
			break
		}
		//doing..
		for _, repo := range repos {
			nameRepos = append(nameRepos, repo.Name)
			total += 1
			fmt.Println("GetRepoList : ", total, repo.Name)
		}
	}
	return nameRepos
}

func (team Team) UpdateRepoPermissionTeam(permission string, teamName string, repoName string) {

	type Payload struct {
		Permission string `json:"permission"`
	}

	data := Payload{
		Permission: permission,
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	body := bytes.NewReader(payloadBytes)

	github := GitHub{Auth: team.Auth}
	statusCode, bodyBytes := github.Request("PUT", "https://api.github.com/orgs/corp-ais/teams/"+teamName+"/repos/"+team.Owner+"/"+repoName, body)

	if statusCode == 204 {
		fmt.Println("SUCCESS : Add/Update PERMISSION [ " + permission + " ] TEAM  [ " + teamName + " ] to REPO NAME [ " + repoName + " ]")
	} else {
		fmt.Println("ERROR : update PERMISSION [ " + permission + " ] TEAM  [ " + teamName + " ] REPO NAME [ " + repoName + " ]")
		log.Println(statusCode, github.GetMessage(bodyBytes))
	}
}

func (team Team) AddnewTeamInAnotherRepoTeam(teamNameAdd string, teamNameIsMember string, permission string) {
	for _, repoName := range team.GetRepoList(teamNameIsMember) {
		fmt.Println(repoName)
		team.UpdateRepoPermissionTeam(permission, teamNameAdd, repoName)
	}
	fmt.Println("Add Team : [ " + teamNameAdd + " ] to Repository team [" + teamNameIsMember + "] is member\nPermission is [ " + permission + " ]")
}

func (team Team) GetInfoTeam(teamName string) model.Team {
	github := GitHub{Auth: team.Auth}
	statusCode, bodyBytes := github.Request("GET", "https://api.github.com/orgs/"+team.Owner+"/teams/"+teamName, nil)

	if statusCode != 200 {
		log.Println(statusCode, github.GetMessage(bodyBytes))
	}

	teamm := model.Team{}
	json.Unmarshal(bodyBytes, &teamm)

	return teamm
}
