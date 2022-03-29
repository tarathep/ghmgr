package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/tarathep/ghmgr/login"
	"github.com/tarathep/ghmgr/model"
)

type Team struct {
	Auth  login.Auth
	Owner string
}

func (team Team) GetRepos(teamName string, page string) (error, model.Repos) {
	github := GitHub{Auth: team.Auth}
	_, statusCode, bodyBytes := github.Request("GET", "https://api.github.com/orgs/"+team.Owner+"/teams/"+teamName+"/repos?page="+page, nil)

	if statusCode != 200 {
		log.Println(statusCode, github.GetMessage(bodyBytes))
	}

	repos := model.Repos{}

	json.Unmarshal(bodyBytes, &repos)

	return nil, repos
}

func (team Team) GetRepoList(teamName string) []string {
	var nameRepos []string
	total := 0

	for i := 0; true; i++ {
		pagex := strconv.Itoa((i + 1))

		_, repos := team.GetRepos(teamName, pagex)

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
	_, statusCode, bodyBytes := github.Request("PUT", "https://api.github.com/orgs/corp-ais/teams/"+teamName+"/repos/"+team.Owner+"/"+repoName, body)

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
	_, statusCode, bodyBytes := github.Request("GET", "https://api.github.com/orgs/"+team.Owner+"/teams/"+teamName, nil)

	if statusCode != 200 {
		log.Println(statusCode, github.GetMessage(bodyBytes))
	}

	teamm := model.Team{}
	json.Unmarshal(bodyBytes, &teamm)

	return teamm
}

// ListTeamMemberPerPage  see more : https://docs.github.com/en/rest/reference/teams#list-team-members
func (team Team) ListTeamMemberPerPage(teamName, page, role string) []model.Members {
	github := GitHub{Auth: team.Auth}
	_, statusCode, bodyBytes := github.Request("GET", "https://api.github.com/orgs/"+team.Owner+"/teams/"+teamName+"/members?page="+page+"&role="+role, nil)

	if statusCode != 200 {
		log.Println(statusCode, github.GetMessage(bodyBytes))
	}

	list_team_member := []model.Members{}
	json.Unmarshal(bodyBytes, &list_team_member)

	return list_team_member
}

// ListTeamMember  see more : https://docs.github.com/en/rest/reference/teams#list-team-members
func (team Team) ListTeamMember(teamName string, role string) []model.Members {
	var listTeamMember []model.Members

	for i := 0; true; i++ {
		page := strconv.Itoa((i + 1))

		list_team_member_perpage := team.ListTeamMemberPerPage(teamName, page, role)

		if len(list_team_member_perpage) == 0 {
			break
		}

		for _, team_member := range list_team_member_perpage {
			listTeamMember = append(listTeamMember, team_member)
		}
	}
	return listTeamMember
}

func RemoveIndex(s []model.Members, index int) []model.Members {
	return append(s[:index], s[index+1:]...)
}

// Optional Exclude IBM Team
func (team Team) ListTeamMemberExcludeTeam(teamName string, teamExcude string, role string) []model.Members {
	var listTeamMember []model.Members
	excludeTeamMember := team.ListTeamMember(teamExcude, "all")

	for _, team_member := range team.ListTeamMember(teamName, role) {
		if !isExist(team_member.Login, excludeTeamMember) {
			listTeamMember = append(listTeamMember, team_member)
		}
	}
	return listTeamMember
}

func isExist(team_member_login string, excludeTeamMember []model.Members) bool {
	for _, exteam := range excludeTeamMember {
		if exteam.Login == team_member_login {
			return true
		}
	}
	return false
}

// https://docs.github.com/en/rest/reference/teams#remove-team-membership-for-a-user
func (team Team) RemoveTeamMembershipForUser(teamname, username string) error {
	github := GitHub{Auth: team.Auth}
	err, statusCode, bodyBytes := github.Request("DELETE", "https://api.github.com/orgs/"+team.Owner+"/teams/"+teamname+"/memberships/"+username, nil)

	if err != nil {
		return err
	}
	if statusCode != 204 {
		return errors.New(github.GetMessage(bodyBytes))
	}
	return nil
}

//https://docs.github.com/en/rest/reference/teams#get-team-membership-for-a-user
func (team Team) GetTeamMembershipForUser(teamname, username string) (error, bool, model.MembershipTeam) {
	github := GitHub{Auth: team.Auth}
	err, statusCode, bodyBytes := github.Request("GET", "https://api.github.com/orgs/"+team.Owner+"/teams/"+teamname+"/memberships/"+username, nil)
	if err != nil {
		return err, false, model.MembershipTeam{}
	}

	membership := model.MembershipTeam{}
	json.Unmarshal(bodyBytes, &membership)

	if statusCode == 200 {
		return nil, true, membership
	} else {
		return nil, false, membership
	}
}

// List teams https://docs.github.com/en/rest/reference/teams#list-teams
func (team Team) ListTeams() (error, []model.Team) {
	github := GitHub{Auth: team.Auth}
	err, statusCode, bodyBytes := github.Request("GET", "https://api.github.com/orgs/"+team.Owner+"/teams", nil)

	if err != nil {
		return err, []model.Team{}
	}

	team_ := []model.Team{}
	json.Unmarshal(bodyBytes, &team_)

	if statusCode != 200 {
		return err, team_
	}

	return nil, team_
}
