package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/tarathep/ghmgr/csv"
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
		color.New(color.FgRed).Println(statusCode, github.GetMessage(bodyBytes))
		os.Exit(1)
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

func (team Team) CheckTeamInORG(teamName string) (error, bool) {
	if err, teams := team.ListTeams(); err != nil {
		return err, false
	} else {
		for _, t := range teams {
			if t.Slug == teamName {
				return nil, true
			}
		}
		return nil, false
	}
}

// List teams https://docs.github.com/en/rest/reference/teams#list-teams
func (team Team) ListTeams() (error, []model.Team) {
	var listTeam []model.Team

	for i := 0; true; i++ {
		page := strconv.Itoa((i + 1))

		err, list_team_perpage := team.ListTeamsPerPage(page)
		if err != nil {
			return err, nil
		}

		if len(list_team_perpage) == 0 {
			break
		}

		for _, teams := range list_team_perpage {
			listTeam = append(listTeam, teams)
		}
	}
	return nil, listTeam
}

func (team Team) ListTeamsPerPage(page string) (error, []model.Team) {
	github := GitHub{Auth: team.Auth}
	err, statusCode, bodyBytes := github.Request("GET", "https://api.github.com/orgs/"+team.Owner+"/teams?page="+page, nil)

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

func (team Team) MembershipOfTeams(username string) (err error, output []model.Team) {
	err, teams := team.ListTeams()
	if err != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}
	for _, team_ := range teams {
		err, isMember, _ := team.GetTeamMembershipForUser(team_.Name, username)

		if err != nil {
			return err, nil
		}
		if isMember {
			output = append(output, team_)
		}
	}
	return nil, output
}

// CHECK USERNAME IN TEAMS CACHE
func (team Team) MembershipOfTeamsCacheTeam(username string) (error, []string) {
	var teams []string
	err, cacheTeams := team.ImportTeamMemberCache()
	if err != nil {
		return err, nil
	}
	for _, cts := range cacheTeams {
		for _, ct := range cts {

			if ct.Username == username {
				teams = append(teams, ct.Team)
			}
		}
	}
	return nil, teams
}

func (team Team) ImportTeamMemberCache() (error, [][]model.Cache) {
	files, err := ioutil.ReadDir("./cache/teams")
	if err != nil {
		return err, nil
	}
	var cacheTeams [][]model.Cache
	for _, f := range files {
		_, c := GetCache("./cache/teams/" + f.Name())
		cacheTeams = append(cacheTeams, c)
	}

	return nil, cacheTeams
}

func (team Team) MembershipOfTeamsCache(caches []model.Cache, username string) []string {
	for _, c := range caches {
		if c.Username == username {
			return strings.Split(c.Team, ",")
		}
	}
	return nil
}

func (team Team) CheckMembershipOutOfTeamsCache(caches []model.Cache, username string) bool {
	for _, c := range caches {
		if c.Username == username {
			if c.Team == "" && c.Username != "" {
				return true
			}
		}
	}
	return false
}

func (team Team) MemberCacheByUser(caches []model.Cache, username string) model.Cache {
	for _, c := range caches {
		if c.Username == username {
			return c
		}
	}
	return model.Cache{}
}

func (team Team) CSVTemplate(templates []model.ProjectMemberListTemplate, email string) model.ProjectMemberListTemplate {
	for _, template := range templates {
		if template.Email == email && template.GitHub == "Y" {
			return template
		}
	}
	return model.ProjectMemberListTemplate{}
}

func SetCache(name string, cache []model.Cache) {
	csv.Template{}.WriteCache(name, cache)
}

func GetCache(name string) (error, []model.Cache) {
	err, models := csv.Template{}.ReadCache(name)

	if err != nil {
		return err, nil
	}

	return nil, models
}

func bstFindByID(users []int, x int) bool {
	i := sort.Search(len(users), func(i int) bool { return x <= users[i] })

	if i < len(users) && users[i] == x {
		return true
	}
	return false
}

func (team Team) AddOrUpdateTeamMembershipForAUser(username string, teamName string, role string) (error, model.TeamRole) {

	payloadBytes, err := json.Marshal(struct {
		Role string `json:"role"`
	}{Role: role})

	if err != nil {
		return err, model.TeamRole{}
	}
	body := bytes.NewReader(payloadBytes)

	github := GitHub{Auth: team.Auth}
	_, statusCode, bodyBytes := github.Request("PUT", "https://api.github.com/orgs/"+team.Owner+"/teams/"+teamName+"/memberships/"+username, body)

	if statusCode != 200 {
		return errors.New(github.GetMessage(bodyBytes)), model.TeamRole{}
	}

	teamRole := model.TeamRole{}
	json.Unmarshal(bodyBytes, &teamRole)

	return nil, teamRole
}
