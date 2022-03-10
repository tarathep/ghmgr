package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tarathep/ghmgr/login"
	"github.com/tarathep/ghmgr/model"
)

type Member struct {
	Auth  login.Auth
	Owner string
}

func (member Member) InviteToCorpTeam(Email string, Role string, teamID int) {

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
		return
	}

	fmt.Println("Invitation Success", Email, Role)

	// fmt.Print(github.GetMessage(bodyBytes))

	// teamm := model.Team{}
	// json.Unmarshal(bodyBytes, &teamm)

	// fmt.Println(teamm)
}

// InvitedToCorpTeamPending : https://docs.github.com/en/rest/reference/teams#list-pending-team-invitations
func (member Member) ListPendingTeamInvitations(teamName string) model.InvitationList {

	github := GitHub{Auth: member.Auth}
	statusCode, bodyBytes := github.Request("GET", "https://api.github.com/orgs/"+member.Owner+"/teams/"+teamName+"/invitations", nil)

	if statusCode != 200 {
		log.Println(statusCode, github.GetMessage(bodyBytes))
		return nil
	}

	invitations := model.InvitationList{}
	json.Unmarshal(bodyBytes, &invitations)

	for i, invitation := range invitations {
		fmt.Println(i, invitation.ID, invitation.Email)
	}
	return invitations
}
