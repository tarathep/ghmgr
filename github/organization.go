package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/tarathep/ghmgr/csv"
	"github.com/tarathep/ghmgr/login"
	"github.com/tarathep/ghmgr/model"
)

type Organization struct {
	Auth  login.Auth
	Owner string
}

//more : https://docs.github.com/en/rest/reference/orgs#create-an-organization-invitation
func (organization Organization) InviteToCorpTeamUserName(username string, role string, teamID int) error {

	// username to user id
	err, usr := User{organization.Auth}.UserInfo(username)
	if err != nil {
		return err
	}

	type Payload struct {
		Invitee_ID int    `json:"invitee_id"`
		Role       string `json:"role"`
		Team_IDs   []int  `json:"team_ids"`
	}
	return organization.createOrganizationInvitation(Payload{
		Invitee_ID: usr.ID,
		Role:       role,
		Team_IDs:   []int{teamID},
	})
}

// https://docs.github.com/en/rest/reference/orgs#check-organization-membership-for-a-user
func (organization Organization) CheckOrganizationMembership(username string) (error, string) {
	github := GitHub{Auth: organization.Auth}
	_, statusCode, bodyBytes := github.Request("GET", "https://api.github.com/orgs/"+organization.Owner+"/members/"+username, nil)

	if statusCode == 204 {
		return nil, "204"
	} else if statusCode == 302 {
		return nil, "302"
	}

	return errors.New(github.GetMessage(bodyBytes)), ""
}

// https://docs.github.com/en/rest/reference/orgs#remove-an-organization-member
func (organization Organization) RemoveOrganizationMember(username string) error {
	github := GitHub{Auth: organization.Auth}
	_, statusCode, bodyBytes := github.Request("DELETE", "https://api.github.com/orgs/"+organization.Owner+"/members/"+username, nil)

	if statusCode == 204 {
		return nil
	}
	return errors.New(github.GetMessage(bodyBytes))
}

// https://docs.github.com/en/rest/reference/orgs#cancel-an-organization-invitation
func (organization Organization) CancelOrganizationInvitation(invitaionID string) error {
	github := GitHub{Auth: organization.Auth}
	_, statusCode, bodyBytes := github.Request("DELETE", "https://api.github.com/orgs/"+organization.Owner+"/invitations/"+invitaionID, nil)

	if statusCode != 204 {
		return errors.New(github.GetMessage(bodyBytes))
	}
	return nil
}

// https://docs.github.com/en/rest/reference/orgs#create-an-organization-invitation
func (organization Organization) createOrganizationInvitation(data interface{}) error {
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	github := GitHub{Auth: organization.Auth}
	_, statusCode, bodyBytes := github.Request("POST", "https://api.github.com/orgs/"+organization.Owner+"/invitations", body)

	if statusCode != 200 && statusCode != 201 {
		return errors.New(github.GetMessage(bodyBytes))
	}
	return nil
}

func (organization Organization) InviteToCorpTeam(email string, role string, teamID int) error {

	if role == "member" || role == "Member" || role == "direct_member" {
		role = "direct_member"
	} else {
		return errors.New("Invalid role " + role + " not support!")
	}

	type Payload struct {
		Email    string `json:"email"`
		Role     string `json:"role"`
		Team_IDs []int  `json:"team_ids"`
	}

	return organization.createOrganizationInvitation(Payload{
		Email:    email,
		Role:     role,
		Team_IDs: []int{teamID},
	})
}

// InvitedToCorpTeamPending : https://docs.github.com/en/rest/reference/teams#list-pending-team-invitations
func (organization Organization) ListPendingTeamInvitations(teamName string) (error, []model.Invitation) {
	var listInvitaion []model.Invitation

	for i := 0; true; i++ {
		page := strconv.Itoa((i + 1))

		err, list_invitation_perpage := organization.ListPendingTeamInvitationsPerPage(teamName, page)
		if err != nil {
			return err, nil
		}
		if len(list_invitation_perpage) == 0 {
			break
		}

		for _, team_member := range list_invitation_perpage {
			listInvitaion = append(listInvitaion, team_member)
		}
	}
	return nil, listInvitaion

}

// https://docs.github.com/en/rest/reference/orgs#list-pending-organization-invitations
func (organization Organization) ListPendingTeamInvitationsPerPage(teamName string, page string) (error, []model.Invitation) {

	github := GitHub{Auth: organization.Auth}
	_, statusCode, bodyBytes := github.Request("GET", "https://api.github.com/orgs/"+organization.Owner+"/teams/"+teamName+"/invitations?per_page=30&page="+page, nil)

	if statusCode != 200 {
		return errors.New(github.GetMessage(bodyBytes)), nil
	}

	invitations := []model.Invitation{}
	json.Unmarshal(bodyBytes, &invitations)

	return nil, invitations
}

// https://docs.github.com/en/rest/reference/teams#list-team-members
func (organization Organization) ListOrgMemberPerPage(page string) (error, []model.Members) {

	github := GitHub{Auth: organization.Auth}
	_, statusCode, bodyBytes := github.Request("GET", "https://api.github.com/orgs/"+organization.Owner+"/members?page="+page, nil)

	if statusCode != 200 {
		return errors.New(github.GetMessage(bodyBytes)), nil
	}

	list_org_member := []model.Members{}
	json.Unmarshal(bodyBytes, &list_org_member)

	return nil, list_org_member
}

//List organization members
func (organization Organization) ListOrgMember() (error, []model.Members) {
	var listOrgMember []model.Members

	for i := 0; true; i++ {
		page := strconv.Itoa((i + 1))

		err, list_org_member_perpage := organization.ListOrgMemberPerPage(page)
		if err != nil {
			return err, nil
		}
		if len(list_org_member_perpage) == 0 {
			break
		}

		for _, team_member := range list_org_member_perpage {
			listOrgMember = append(listOrgMember, team_member)
		}
	}
	return nil, listOrgMember
}

func (organization Organization) SetCache(cache []model.Cache) {
	csv.Template{}.WriteCache(cache)
}

func (organization Organization) GetCache() (error, []model.Cache) {
	err, models := csv.Template{}.ReadCache()

	if err != nil {
		return err, nil
	}

	// for _, d := range models {
	// 	fmt.Println(d.ID)
	// }
	return nil, models
}
