package manage

import (
	"fmt"
	"log"

	"github.com/tarathep/ghmgr/csv"
	"github.com/tarathep/ghmgr/github"
	"github.com/tarathep/ghmgr/model"
)

type GitHubManager struct {
	github.Team
	github.Member
	github.User
}

func (mgr GitHubManager) ShowListTeamMember(teamName string, role string, email string) {
	fmt.Println("List Team Member of " + teamName)

	for i, teamMember := range mgr.Team.ListTeamMember(teamName, role) {
		if !(role == "all" || role == "member" || role == "maintainer") {
			log.Fatal("invalid role")
		} else if role == "all" {
			if email == "show" {
				_, usrInfo := mgr.UserInfo(teamMember.Login)
				println(i+1, "ID:", teamMember.ID, "User:", teamMember.Login, "Email:", usrInfo.Email)
			} else {
				println(i+1, "ID:", teamMember.ID, "User:", teamMember.Login)
			}
		} else {
			println(i+1, "ID:", teamMember.ID, "User:", teamMember.Login, "Role:", role)
		}
	}
}

func (mgr GitHubManager) ShowListTeamMemberExclude(teamName string, teamExcude string, role string, email string) {
	fmt.Println("List Team Member of " + teamName)

	for i, teamMember := range mgr.Team.ListTeamMemberExcludeTeam(teamName, teamExcude, role) {
		if !(role == "all" || role == "member" || role == "maintainer") {
			log.Fatal("invalid role")
		} else if role == "all" {
			if email == "show" {
				_, usrInfo := mgr.UserInfo(teamMember.Login)
				println(i+1, "ID:", teamMember.ID, "User:", teamMember.Login, "Email:", usrInfo.Email)
			} else {
				println(i+1, "ID:", teamMember.ID, "User:", teamMember.Login)
			}
		} else {
			println(i+1, "ID:", teamMember.ID, "User:", teamMember.Login, "Role:", role)
		}
	}
}

func (mgr GitHubManager) ReadCSVFile(fileName string) {
	templ := csv.Template{}

	proj, csvTemplate := templ.ReadFile(fileName)

	for i, csvTempl := range csvTemplate {
		fmt.Println((i + 1), proj, csvTempl.Email)
	}
}

func (mgr GitHubManager) InviteMemberToCorpTeam(teamName string, role string, email string) {

	if role == "member" || role == "Member" {
		role = "direct_member"
	} else if role == "maintainer" || role == "Maintainer" {
		role = "maintainer"
	}

	teamID := mgr.Team.GetInfoTeam(teamName).ID
	if err := mgr.Member.InviteToCorpTeam(email, role, teamID); err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Print("Invite Successful")
	}
}

func (mgr GitHubManager) InviteMemberToCorpTeamUsername(teamName string, role string, username string) {

	if role == "member" || role == "Member" {
		role = "direct_member"
	} else if role == "maintainer" || role == "Maintainer" {
		role = "maintainer"
	}

	teamID := mgr.Team.GetInfoTeam(teamName).ID
	if err := mgr.Member.InviteToCorpTeamUserName(username, role, teamID); err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Print("Invite Successful")
	}
}

func (mgr GitHubManager) InviteMemberToCorpTeamCSV(fileName string) {

	templ := csv.Template{}

	proj, csvTemplate := templ.ReadFile(fileName)

	for i, csvTempl := range csvTemplate {
		fmt.Println((i + 1), proj, csvTempl.Email)

		mgr.InviteMemberToCorpTeam(proj, "direct_member", csvTempl.Email)
	}
}

func (mgr GitHubManager) ShowListTeamMemberPending(teamName string) {
	fmt.Println("List Team Member of " + teamName + "\n------------------------------------------------")

	for i, invitation := range mgr.Member.ListPendingTeamInvitations(teamName) {
		println(i+1, "ID:", invitation.ID, "Email:", invitation.Email)
	}
}

func (mgr GitHubManager) ExportCSVMemberTeam(teamName string) {
	var dataset []model.CSV

	for _, role := range []string{"maintainer", "member"} {
		for _, teamMember := range mgr.Team.ListTeamMember(teamName, role) {
			dataset = append(dataset, model.CSV{GitHubUser: teamMember.Login, GitHubTeamRole: role})
		}
	}
	result := csv.Template{}.WriteCSV(teamName, dataset)

	if result {
		fmt.Print("Export CSV Member Team : " + teamName)
	}
}

func (mgr GitHubManager) ExportCSVMemberTeamExclude(teamName string, teamExclude string) {
	var dataset []model.CSV

	for _, role := range []string{"maintainer", "member"} {
		for _, teamMember := range mgr.Team.ListTeamMemberExcludeTeam(teamName, teamExclude, role) {
			dataset = append(dataset, model.CSV{GitHubUser: teamMember.Login, GitHubTeamRole: role})
		}
	}
	result := csv.Template{}.WriteCSV(teamName, dataset)

	if result {
		fmt.Print("Export CSV Member Team : " + teamName)
	}
}

func (mgr GitHubManager) CancelOrganizationInvitation(username string) {
	if err := mgr.Member.CancelOrganizationInvitation(username); err != nil {
		log.Fatal("cancel invite err :", username, err)
	} else {
		log.Print("Cancel invite", username, "success")
	}
}

func (mgr GitHubManager) CheckOrganizationMembership(username string) {
	if err, body := mgr.Member.CheckOrganizationMembership(username); err != nil {
		log.Print(err.Error())
	} else if body == "204" {
		fmt.Print(username, " is an organization member and user is a member")
	} else if body == "302" {
		fmt.Print(username, " is not an organization member")
	}
}
