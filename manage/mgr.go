package manage

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
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

	color.New(color.Italic).Print("Team members will include the members of child teams.\nTo list members in a [" + teamName + "] team, the team must be visible to the authenticated user.\n")
	if email != "" {
		color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%23s\t%40s\n", "No.", "ID", "Username", "Email")
	} else {
		color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%23s\n", "No.", "ID", "Username")
	}
	for i, teamMember := range mgr.Team.ListTeamMember(teamName, role) {
		if !(role == "all" || role == "member" || role == "maintainer") {
			color.New(color.FgRed).Println("Invalid role")
			os.Exit(1)
		} else if role == "all" {
			if email == "show" {
				_, usrInfo := mgr.UserInfo(teamMember.Login)
				fmt.Printf("%3d\t%10d\t%23s\t%40s\n", i+1, teamMember.ID, teamMember.Login, usrInfo.Email)
			} else {
				fmt.Printf("%3d\t%10d\t%23s\n", i+1, teamMember.ID, teamMember.Login)
			}
		} else {
			fmt.Printf("%3d\t%10d\t%23s\n", i+1, teamMember.ID, teamMember.Login)
		}
	}
}

func (mgr GitHubManager) ShowListTeamMemberExclude(teamName string, teamExcude string, role string, email string) {

	color.New(color.Italic).Print("Team members will include the members of child teams.\nTo list members in a [" + teamName + "] team and Exclude [" + teamExcude + "] , the team must be visible to the authenticated user.\n")

	if email != "" {
		color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%23s\t%40s\n", "No.", "ID", "Username", "Email")
	} else {
		color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%23s\n", "No.", "ID", "Username")
	}
	for i, teamMember := range mgr.Team.ListTeamMemberExcludeTeam(teamName, teamExcude, role) {
		if !(role == "all" || role == "member" || role == "maintainer") {
			color.New(color.FgRed).Println("Invalid role")
			os.Exit(1)
		} else if role == "all" {
			if email == "show" {
				_, usrInfo := mgr.UserInfo(teamMember.Login)
				fmt.Printf("%3d\t%10d\t%23s\t%40s\n", i+1, teamMember.ID, teamMember.Login, usrInfo.Email)
			} else {
				fmt.Printf("%3d\t%10d\t%23s\n", i+1, teamMember.ID, teamMember.Login)
			}
		} else {
			fmt.Printf("%3d\t%10d\t%23s\n", i+1, teamMember.ID, teamMember.Login)
		}
	}
}

func (mgr GitHubManager) ReadCSVFile(fileName string) {

	templ := csv.Template{}

	err, proj, csvTemplate := templ.ReadFile(fileName)
	if err != nil {
		color.New(color.FgHiRed).Println(err.Error())
		os.Exit(1)
	}

	color.New(color.Italic).Print("CSV File Reader.\nTo list members in a  CSV file , [" + proj + "] team, the team must be visible to the GitHub.\n")

	color.New(color.FgHiMagenta).Printf("%2s\t%5s\t%30s\t%40s\t%10s\t%10s\t%15s\n", "No.", "ID", "MemberName", "Email", "Role", "Team Role", "UserName")
	// fmt.Printf("%s\n", strings.Repeat("-", 170-27))
	for i, csvTempl := range csvTemplate {
		fmt.Printf("%2d\t%5s\t%30s\t%40s\t%10s\t%10s\t%15s\n", (i + 1), csvTempl.ID, csvTempl.MemberName, csvTempl.Email, csvTempl.Role, csvTempl.GitHubTeamRole, csvTempl.GitHubUser)

	}
}

func (mgr GitHubManager) InviteMemberToCorpTeamEmail(teamName string, role string, email string) {
	color.New(color.Italic).Print("Create an organization invitation assign to [" + teamName + "] team. \n")
	mgr.InviteMemberToCorpTeam(teamName, role, email)
}

func (mgr GitHubManager) InviteMemberToCorpTeam(teamName string, role string, email string) {
	fmt.Printf(" %40s\t%20s : ", email, teamName)

	//  UNCOMMENT to INVITE !!!!

	// fmt.Printf("%3d\t%10d\t%40s\n", i+1, invitation.ID, invitation.Email)

	// teamID := mgr.Team.GetInfoTeam(teamName).ID

	// if err := mgr.Member.InviteToCorpTeam(email, role, teamID); err != nil {
	// 	color.New(color.FgHiRed).Println("Error ", err.Error())
	// 	os.Exit(1)
	// } else {

	// }
	color.New(color.FgHiGreen).Println("Done")
}

//  xxx reflactor waiting
func (mgr GitHubManager) InviteMemberToCorpTeamUsername(teamName string, role string, username string) {

	// if role == "member" || role == "Member" {
	// 	role = "direct_member"
	// } else if role == "maintainer" || role == "Maintainer" {
	// 	role = "maintainer"
	// }

	teamID := mgr.Team.GetInfoTeam(teamName).ID
	if err := mgr.Member.InviteToCorpTeamUserName(username, role, teamID); err != nil {
		log.Fatal(err.Error())

		os.Exit(1)
	} else {
		fmt.Print("Invite Successful")
	}
}

func (mgr GitHubManager) InviteMemberToCorpTeamCSV(fileName string) {

	color.New(color.Italic).Print("Create an organization invitation from [" + fileName + "] file. \n")

	templ := csv.Template{}

	err, proj, csvTemplate := templ.ReadFile(fileName)
	if err != nil {
		color.New(color.FgHiRed).Println(err.Error())
		os.Exit(1)
	}

	for i, csvTempl := range csvTemplate {
		fmt.Print((i + 1), "\t")
		mgr.InviteMemberToCorpTeam(proj, "direct_member", csvTempl.Email)
	}
}

func (mgr GitHubManager) ShowListTeamMemberPending(teamName string) {
	color.New(color.Italic).Print("List pending [" + teamName + "] team invitations\n")

	err, pendings := mgr.Member.ListPendingTeamInvitations(teamName)
	if err != nil {
		color.New(color.FgHiRed).Println(err.Error())
		os.Exit(1)
	}

	color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%40s\n", "No.", "ID", "Email")

	for i, invitation := range pendings {
		fmt.Printf("%3d\t%10d\t%40s\n", i+1, invitation.ID, invitation.Email)
	}
}

func (mgr GitHubManager) ExportCSVMemberTeam(teamName string) {

	color.New(color.Italic).Print("Export CSV Member Team [" + teamName + "] : ")
	var dataset []model.CSV

	for _, role := range []string{"maintainer", "member"} {
		for _, teamMember := range mgr.Team.ListTeamMember(teamName, role) {
			dataset = append(dataset, model.CSV{GitHubUser: teamMember.Login, GitHubTeamRole: role})
		}
	}
	result := csv.Template{}.WriteCSV(teamName, dataset)

	if result {
		color.New(color.FgHiGreen).Print("Done")
	}
}

func (mgr GitHubManager) ExportCSVMemberTeamExclude(teamName string, teamExclude string) {
	color.New(color.Italic).Print("Export CSV Member Team [" + teamName + "] Exclude Team [" + teamExclude + "] : ")
	var dataset []model.CSV

	for _, role := range []string{"maintainer", "member"} {
		for _, teamMember := range mgr.Team.ListTeamMemberExcludeTeam(teamName, teamExclude, role) {
			dataset = append(dataset, model.CSV{GitHubUser: teamMember.Login, GitHubTeamRole: role})
		}
	}
	result := csv.Template{}.WriteCSV(teamName, dataset)

	if result {
		color.New(color.FgHiGreen).Print("Done")
	}
}

func (mgr GitHubManager) CancelOrganizationInvitation(username string) {
	color.New(color.Italic).Println("Cancel an organization invitation. In order to cancel an organization invitation, the authenticated user must be an organization owner.")

	if err := mgr.Member.CancelOrganizationInvitation(username); err != nil {
		color.New(color.FgHiRed).Print("cancel invite err :", username, err)
		os.Exit(1)
	} else {
		color.New(color.FgYellow).Print("Cancel invite", username, "success")
	}
}

func (mgr GitHubManager) CheckOrganizationMembership(username string) {
	color.New(color.Italic).Println("Check if a user is, publicly or privately, a member of the organization.")

	if err, _ := mgr.Member.CheckOrganizationMembership(username); err == nil {
		color.New(color.FgHiGreen).Print(username, " is an organization member and user is a member")
	} else {
		color.New(color.FgHiRed).Print(username, " is not an organization member or err ", err.Error())
	}
}

func (mgr GitHubManager) RemoveOrganizationMember(username string) {

	color.New(color.Italic).Print("Removing a user from this list will remove them from all teams and they will no longer have any access to the organization's repositories\n")
	if err := mgr.Member.RemoveOrganizationMember(username); err != nil {
		log.Fatal(err.Error())

	}
	color.New(color.FgHiRed).Print(username, " was removed an organization")
}

func (mgr GitHubManager) RemoveTeamMembershipForUser(teamname, username string) {
	color.New(color.Italic).Print("To remove a membership between a user and a team, the authenticated user must have 'admin' permissions to the team or be an owner of the organization that the team is associated with. Removing team membership does not delete the user, it just removes their membership from the team.\n")
	color.New(color.FgHiRed).Print(username, " removing a "+teamname+" team :")
	color.New(color.FgHiGreen).Println(" Done")

}
