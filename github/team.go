package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/tarathep/githuby/login"
	"github.com/tarathep/githuby/model"
)

type Team struct {
	Auth  login.Auth
	Debug bool
}

func (team Team) GetRepos(owner string, teamName string, page string) model.Repos {
	req, err := http.NewRequest("GET", "https://api.github.com/orgs/"+owner+"/teams/"+teamName+"/repos?page="+page, nil)
	if err != nil {
		// handle err
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	//req.Header.Set("Authorization", "token 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")
	req.Header.Set("Authorization", "token "+team.Auth.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// bodyString := string(bodyBytes)

	// fmt.Println(bodyString)

	repos := model.Repos{}

	json.Unmarshal(bodyBytes, &repos)

	// log.Print(repos)

	return repos
}

func (team Team) GetRepoList(owner string, teamName string) []string {
	var nameRepos []string
	total := 0

	for i := 0; true; i++ {
		page := strconv.Itoa((i + 1))
		repos := team.GetRepos(owner, teamName, page)

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

func (team Team) UpdateRepoPermissionTeam(permission string, teamName string, owner string, repoName string) {

	type Payload struct {
		Permission string `json:"permission"`
	}

	data := Payload{
		Permission: permission,
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("PUT", "https://api.github.com/orgs/corp-ais/teams/"+teamName+"/repos/"+owner+"/"+repoName, body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authorization", "token "+team.Auth.Token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		fmt.Println("SUCCESS : Add/Update PERMISSION [ " + permission + " ] TEAM  [ " + teamName + " ] to REPO NAME [ " + repoName + " ]")
	} else {
		fmt.Println("ERROR : update PERMISSION [ " + permission + " ] TEAM  [ " + teamName + " ] REPO NAME [ " + repoName + " ]")

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}
}

func (team Team) AddnewTeamInAnotherRepoTeam(owner string, teamNameAdd string, teamNameIsMember string, permission string) {
	for _, repoName := range team.GetRepoList(owner, teamNameIsMember) {
		fmt.Println(repoName)
		team.UpdateRepoPermissionTeam(permission, teamNameAdd, owner, repoName)
	}
	fmt.Println("Add Team : [ " + teamNameAdd + " ] to Repository team [" + teamNameIsMember + "] is member\nPermission is [ " + permission + " ]")
}

func (team Team) GetInfoTeam(owner string, teamName string) model.Team {
	req, err := http.NewRequest("GET", "https://api.github.com/orgs/"+owner+"/teams/"+teamName, nil)
	if err != nil {
		// handle err
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "token "+team.Auth.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		bodyString := string(bodyBytes)
		log.Fatal(resp.Status, bodyString)
	}

	teamm := model.Team{}
	json.Unmarshal(bodyBytes, &teamm)

	return teamm
}
