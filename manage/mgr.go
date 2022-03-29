package manage

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/tarathep/ghmgr/csv"
	"github.com/tarathep/ghmgr/github"
	"github.com/tarathep/ghmgr/model"
)

type GitHubManager struct {
	github.Team
	github.Organization
	github.User
}

// https://docs.github.com/en/rest/reference/teams#get-team-membership-for-a-user
func (mgr GitHubManager) CheckTeamMembershipForUser(teamName string, username string) {
	color.New(color.Italic).Print("Get team membership for a user\nTeam members will include the members of child teams.\nTo get a user's membership with a team, the team must be visible to the authenticated user.\n")

	err, isMember, membership := mgr.Team.GetTeamMembershipForUser(teamName, username)

	if err != nil {
		fmt.Println(err)
		return
	}
	if isMember {
		fmt.Println(isMember)
	}
	fmt.Println(membership.Role)

	// if err, is, membership := mgr.Team.GetTeamMembershipForUser(teamName, username); err == nil && !(is) {
	// 	color.New(color.FgHiGreen).Print(username, " is a member of team "+teamName+" and Role : "+membership.Role)
	// } else {
	// 	color.New(color.FgHiRed).Print(username, " is not a team member "+teamName+" or err ", err.Error())
	// }
}

// List teams https://docs.github.com/en/rest/reference/teams#list-teams
func (mgr GitHubManager) ListTeam() {
	color.New(color.Italic).Print("Lists all teams in an organization that are visible to the authenticated user.\n")

	color.New(color.FgHiMagenta).Println("No.", "\tID", "\t\tTeam Name")

	err, teams := mgr.Team.ListTeams()
	if err != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}

	for i, team := range teams {
		fmt.Println(i+1, "\t"+strconv.Itoa(team.ID), "\t"+team.Name)
	}
}

//https://docs.github.com/en/rest/reference/teams#list-team-members
func (mgr GitHubManager) ListTeamMembers(option string) {
	color.New(color.Italic).Print("Team members will include the members of child teams.\nTo list members in a team, the team must be visible to the authenticated user..\n")

	if option == "email" {
		color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%23s\t%40s\n", "No.", "ID", "Username", "Email")
	} else {
		color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%23s\n", "No.", "ID", "Username")
	}

	err, i := mgr.Organization.ListOrgMember()
	if err != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}

	err, load := mgr.GetCache()
	if err == nil {
		// CHCHE
		for i, l := range load {
			if option == "email" {
				fmt.Printf("%3d\t%10s\t%23s\t%40s\n", i+1, l.ID, l.Username, l.Email)

			} else {
				fmt.Printf("%3d\t%10s\t%23s\n", i+1, l.ID, l.Username)
			}
		}
	} else {
		//ORIGINAL
		var cache []model.Cache

		for i, member := range i {

			if option == "email" {
				_, usrInfo := mgr.UserInfo(member.Login)
				fmt.Printf("%3d\t%10d\t%23s\t%40s\n", i+1, member.ID, member.Login, usrInfo.Email)
				cache = append(cache, model.Cache{
					No:       strconv.Itoa(i + 1),
					ID:       strconv.Itoa(member.ID),
					Username: member.Login,
					Email:    usrInfo.Email,
					Team:     "",
				})
			} else {
				fmt.Printf("%3d\t%10d\t%23s\n", i+1, member.ID, member.Login)

				cache = append(cache, model.Cache{
					No:       strconv.Itoa(i + 1),
					ID:       strconv.Itoa(member.ID),
					Username: member.Login,
					Email:    "",
					Team:     "",
				})
			}
		}
		mgr.SetCache(cache)
	}

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
	color.New(color.Italic).Print("Create an organization invitation assign to [" + teamName + "] team. (org support member only) \n")

	//MEMBER ONLY!!
	role = "direct_member"
	mgr.InviteMemberToCorpTeam(teamName, role, email)
}

func (mgr GitHubManager) InviteMemberToCorpTeam(teamName string, role string, email string) {
	fmt.Printf(" %40s\t%20s : ", email, teamName)

	//  UNCOMMENT to INVITE !!!!

	teamID := mgr.Team.GetInfoTeam(teamName).ID

	if err := mgr.Organization.InviteToCorpTeam(email, role, teamID); err != nil {
		color.New(color.FgHiRed).Println("Error ", err.Error())
		os.Exit(1)
	} else {
		color.New(color.FgHiGreen).Println("Done")
	}

}

//  xxx reflactor waiting
func (mgr GitHubManager) InviteMemberToCorpTeamUsername(teamName string, role string, username string) {

	// if role == "member" || role == "Member" {
	// 	role = "direct_member"
	// } else if role == "maintainer" || role == "Maintainer" {
	// 	role = "maintainer"
	// }

	teamID := mgr.Team.GetInfoTeam(teamName).ID
	if err := mgr.Organization.InviteToCorpTeamUserName(username, role, teamID); err != nil {
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

	err, pendings := mgr.Organization.ListPendingTeamInvitations(teamName)
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

// https://docs.github.com/en/rest/reference/orgs#cancel-an-organization-invitation
func (mgr GitHubManager) CancelOrganizationInvitation(invitationID string) {
	color.New(color.Italic).Println("Cancel an organization invitation. In order to cancel an organization invitation, the authenticated user must be an organization owner.")

	color.New(color.FgYellow).Print("Cancel invitation ID [" + invitationID + "] : ")
	if err := mgr.Organization.CancelOrganizationInvitation(invitationID); err != nil {
		color.New(color.FgHiRed).Println("ERROR ", err)
		os.Exit(1)
	} else {
		color.New(color.FgHiGreen).Println("Done")
	}
}

func (mgr GitHubManager) CheckOrganizationMembership(username string) {
	color.New(color.Italic).Println("Check if a user is, publicly or privately, a member of the organization.")

	if err, _ := mgr.Organization.CheckOrganizationMembership(username); err == nil {
		color.New(color.FgHiGreen).Print(username, " is an organization member and user is a member")
	} else {
		color.New(color.FgHiRed).Print(username, " is not an organization member or err ", err.Error())
	}
}

func (mgr GitHubManager) RemoveOrganizationMember(username string) {

	color.New(color.Italic).Print("Removing a user from this list will remove them from all teams and they will no longer have any access to the organization's repositories\n")
	if err := mgr.Organization.RemoveOrganizationMember(username); err != nil {
		log.Fatal(err.Error())

	}
	color.New(color.FgHiRed).Print(username, " was removed an organization")
}

func (mgr GitHubManager) RemoveTeamMembershipForUser(teamname, username string) {
	color.New(color.Italic).Print("To remove a membership between a user and a team, the authenticated user must have 'admin' permissions to the team or be an owner of the organization that the team is associated with. Removing team membership does not delete the user, it just removes their membership from the team.\n")
	color.New(color.FgHiRed).Print(username, " removing a "+teamname+" team :")
	color.New(color.FgHiGreen).Println(" Done")

}
