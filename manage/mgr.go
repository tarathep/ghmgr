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
}

func (mgr GitHubManager) ShowListTeamMember(teamName string, role string) {
	fmt.Println("List Team Member of " + teamName)

	for i, teamMember := range mgr.Team.ListTeamMember(teamName, role) {
		if !(role == "all" || role == "member" || role == "maintainer") {
			log.Fatal("invalid role")
		} else if role == "all" {
			println(i+1, "ID:", teamMember.ID, "User:", teamMember.Login)
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
	// HARDCODE  ROLE!!
	role = "direct_member"
	teamID := mgr.Team.GetInfoTeam(teamName).ID
	mgr.Member.InviteToCorpTeam(email, role, teamID)
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
	fmt.Println("List Team Member of " + teamName)
	for i, invitation := range mgr.Member.ListPendingTeamInvitations(teamName) {
		println(i+1, "ID:", invitation.ID, "User:", invitation.Login, "Email:", invitation.Email)
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
