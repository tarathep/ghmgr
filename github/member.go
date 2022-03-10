package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/tarathep/ghmgr/login"
	"github.com/tarathep/ghmgr/model"
)

type Member struct {
	Auth  login.Auth
	Owner string
}

func (member Member) InviteToCorpTeam(Email string, Role string, teamID int) error {

	type Payload struct {
		// Invitee_ID int    `json:"invitee_id"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		Team_IDs []int  `json:"team_ids"`
	}

	data := Payload{
		Email:    Email,
		Role:     Role,
		Team_IDs: []int{teamID},
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	body := bytes.NewReader(payloadBytes)

	github := GitHub{Auth: member.Auth}
	statusCode, bodyBytes := github.Request("POST", "https://api.github.com/orgs/"+member.Owner+"/invitations", body)

	if statusCode != 200 && statusCode != 201 {
		log.Println(statusCode, github.GetMessage(bodyBytes))
		return errors.New(github.GetMessage(bodyBytes))
	}
	return nil
}

// InvitedToCorpTeamPending : https://docs.github.com/en/rest/reference/teams#list-pending-team-invitations
func (member Member) ListPendingTeamInvitations(teamName string) []model.Invitation {
	var listInvitaion []model.Invitation

	for i := 0; true; i++ {
		page := strconv.Itoa((i + 1))

		list_invitation_perpage := member.ListPendingTeamInvitationsPerPage(teamName, page)

		if len(list_invitation_perpage) == 0 {
			break
		}

		for _, team_member := range list_invitation_perpage {
			listInvitaion = append(listInvitaion, team_member)
		}
	}
	return listInvitaion

}

func (member Member) ListPendingTeamInvitationsPerPage(teamName string, page string) []model.Invitation {

	github := GitHub{Auth: member.Auth}
	statusCode, bodyBytes := github.Request("GET", "https://api.github.com/orgs/"+member.Owner+"/teams/"+teamName+"/invitations?per_page=30&page="+page, nil)

	if statusCode != 200 {
		log.Println(statusCode, github.GetMessage(bodyBytes))
		return nil
	}

	invitations := []model.Invitation{}
	json.Unmarshal(bodyBytes, &invitations)

	return invitations
}

// ListTeamMemberPerPage  see more : https://docs.github.com/en/rest/reference/teams#list-team-members
func (team Team) ListTeamMemberPerPagex(teamName, page, role string) []model.TeamMember {
	github := GitHub{Auth: team.Auth}
	statusCode, bodyBytes := github.Request("GET", "https://api.github.com/orgs/"+team.Owner+"/teams/"+teamName+"/members?page="+page+"&role="+role, nil)

	if statusCode != 200 {
		log.Println(statusCode, github.GetMessage(bodyBytes))
	}

	list_team_member := []model.TeamMember{}
	json.Unmarshal(bodyBytes, &list_team_member)

	return list_team_member
}

func (member Member) ListMember() {

	github := GitHub{Auth: member.Auth}
	statusCode, bodyBytes := github.Request("GET", "https://api.github.com/orgs/"+member.Owner+"/members", nil)

	if statusCode != 200 {
		log.Println(statusCode, github.GetMessage(bodyBytes))
		return
	}

	log.Println(statusCode, github.GetMessage(bodyBytes))

}
