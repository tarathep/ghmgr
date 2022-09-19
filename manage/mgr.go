package manage

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/tarathep/ghmgr/csv"
	"github.com/tarathep/ghmgr/github"
	"github.com/tarathep/ghmgr/model"
)

type GitHubManager struct {
	github.Team
	github.Organization
	github.User
	Version string
}

func LogError(raw interface{}, enable bool) {
	if enable {
		color.New(color.FgHiRed).Print(raw)
	}
}
func LogSuccess(raw interface{}, enable bool) {
	if enable {
		color.New(color.FgHiGreen).Print(raw)
	}
}
func LogInfo(raw interface{}, enable bool) {
	if enable {
		color.New(color.FgHiGreen).Print(raw)
	}
}
func LogWarning(raw interface{}, enable bool) {
	if enable {
		color.New(color.FgHiYellow).Print(raw)
	}
}
func LogDisable(raw interface{}, enable bool) {
	if enable {
		color.New(color.FgHiBlue).Print(raw)
	}
}

func LogCustom(style color.Attribute, raw interface{}, enable bool) {
	if enable {
		color.New(style).Print(raw)
	}
}
func LogPrint(enable bool, a interface{}) {
	if enable {
		fmt.Print(a)
	}
}

func (mgr GitHubManager) CheckTemplateFormat(logging bool, fileName string) bool {
	templ := csv.Template{}
	err, proj, csvTemplate := templ.ReadProjectMemberListTemplateCSV("reports/input/" + fileName)
	if err != nil {
		LogError(err.Error(), logging)
		os.Exit(1)
	}

	LogCustom(color.Italic, "Validating format "+proj+"csv template \n", logging)

	result := true

	for _, csvTempl := range csvTemplate {
		LogPrint(logging, csvTempl.No+"\n")

		if csvTempl.Fullname != "" && regexp.MustCompile(`^[a-zA-Z0-9. ]+$`).MatchString(csvTempl.Fullname) {
			LogSuccess("(✓)", logging)
		} else {
			LogError("(X)", logging)
			result = false
		}
		LogPrint(logging, " Name:"+csvTempl.Fullname+"\n")

		if csvTempl.Email != "" && regexp.MustCompile(`^[a-zA-Z0-9@.]+$`).MatchString(csvTempl.Email) && !(regexp.MustCompile(`@hotmail|@gmail|@outlook|@live|@windowslive|@yahoo`).MatchString(csvTempl.Email)) {
			LogSuccess("(✓)", logging)
		} else {
			LogError("(X)", logging)
			result = false
		}
		LogPrint(logging, " Email:"+csvTempl.Email+"\n")

		if csvTempl.GitHub == "Y" {
			if csvTempl.GitHubUsername != "" {
				if err, _ := mgr.Organization.CheckOrganizationMembership(csvTempl.GitHubUsername); err == nil {
					LogSuccess("(✓)", logging)
				} else {
					LogError("(X)", logging)
					result = false
				}
			} else {
				LogWarning("(!)", logging)
			}
		} else {
			LogDisable("(-)", logging)
		}
		LogPrint(logging, " GitHub Username:"+csvTempl.GitHubUsername+"\n")

		if csvTempl.GitHub == "Y" {
			if !(csvTempl.GitHubTeamRole == "maintainer" && csvTempl.Role == "member") {
				LogSuccess("(✓)", logging)
			} else {
				LogError("(X)", logging)
				result = false
			}
		} else {
			LogDisable("(-)", logging)

		}
		LogPrint(logging, " GitHub Role:"+csvTempl.GitHubTeamRole+"\n")

		if !(csvTempl.GitHub == "Y" && csvTempl.GitHub == "N") {
			LogSuccess("(✓)", logging)
		} else {
			LogError("(X)", logging)
			result = false
		}
		LogPrint(logging, " GitHub :"+csvTempl.GitHub+"\n")
	}

	color.New(color.Italic).Print("Format " + proj + ".csv have been validated : ")
	if result {
		color.New(color.FgHiGreen).Println("(✓)")
		return true
	}
	color.New(color.FgHiRed).Println("(X)")
	return false
}

func (mgr GitHubManager) RewriteTemplateFormat(teamName string) {
	templ := csv.Template{}
	err, proj, csvTemplate := templ.ReadProjectMemberListTemplateCSV("reports/input/" + teamName + ".csv")
	if err != nil {
		color.New(color.FgHiRed).Print(err.Error())
		os.Exit(1)
	}

	color.New(color.Italic).Print("Rewriting format " + proj + ".csv template : ")

	I := 0
	var dataset []model.ProjectMemberListTemplate

	for _, csvTempl := range csvTemplate {
		if csvTempl.Email != "verify email" {
			I++

			if csvTempl.Username == "" && csvTempl.Email != "" {
				csvTempl.Username = strings.Split(csvTempl.Email, "@")[0]
			}
			if csvTempl.SubscriptionOwner == "" {
				csvTempl.SubscriptionOwner = "-"
			}
			if err, _ := mgr.Organization.CheckOrganizationMembership(csvTempl.GitHubUsername); err != nil {
				csvTempl.GitHubUsername = ""
			}
			if csvTempl.GitHubTeamRole != "maintainer" && csvTempl.GitHub == "Y" {
				csvTempl.GitHubTeamRole = "member"
			}
			if csvTempl.GitHub != "Y" {
				csvTempl.GitHub = "N"
			}
			if csvTempl.AzureDEV != "Y" {
				csvTempl.AzureDEV = "N"
			}
			if csvTempl.AzurePRD != "Y" {
				csvTempl.AzurePRD = "N"
			}
			if csvTempl.ELK != "Y" {
				csvTempl.ELK = "N"
			}
			if csvTempl.Jumphost != "Y" {
				csvTempl.Jumphost = "N"
			}
			if csvTempl.Bastion != "Y" {
				csvTempl.Bastion = "N"
			}

			dataset = append(dataset, model.ProjectMemberListTemplate{
				No:                strconv.Itoa(I),
				Username:          csvTempl.Username,
				Fullname:          csvTempl.Fullname,
				Email:             strings.TrimSpace(csvTempl.Email),
				Role:              csvTempl.Role,
				SubscriptionOwner: csvTempl.SubscriptionOwner,
				GitHubUsername:    csvTempl.GitHubUsername,
				GitHubTeamRole:    csvTempl.GitHubTeamRole,
				GitHub:            csvTempl.GitHub,
				AzureDEV:          csvTempl.AzureDEV,
				AzurePRD:          csvTempl.AzurePRD,
				ELK:               csvTempl.ELK,
				Jumphost:          csvTempl.Jumphost,
				Bastion:           csvTempl.Bastion,
			})
		}
	}

	result := csv.Template{}.WriteProjectMemberListTemplateCSV(teamName, "template membership of team Generated by GHMGR "+mgr.Version+" : ", "reports/input/"+teamName, dataset)

	if result != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}
	color.New(color.FgHiGreen).Print("Done\n")
}

func (mgr GitHubManager) GetUsernameFromEmail(email string) {
	color.New(color.Italic).Print("Get Username from email\n")

	if username := mgr.User.EmailToUsername(mgr.loadCache(), email); username != "" {
		color.New(color.FgCyan).Println(username)
		return
	}
	color.New(color.FgYellow).Println("Not Found")
}

func (mgr GitHubManager) GetEmailFromUsername(username string) {
	color.New(color.Italic).Print("Get Email from Username\n")
	_, usr := mgr.UserInfo(username)

	if usr.Email != "" {
		color.New(color.FgCyan).Println(usr.Email)
		return
	}
	color.New(color.FgYellow).Println("Not Found")
}

// https://docs.github.com/en/rest/reference/teams#get-team-membership-for-a-user
func (mgr GitHubManager) CheckTeamMembershipForUser(teamName string, username string) {
	color.New(color.Italic).Print("Get team membership for a user\nTeam members will include the members of child teams.\nTo get a user's membership with a team, the team must be visible to the authenticated user.\n")

	err, isMember, membership := mgr.Team.GetTeamMembershipForUser(teamName, username)

	if err != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}

	if isMember {
		color.New(color.FgHiGreen).Print(username, " is a member of team [ "+teamName+" ] and role [ "+membership.Role+" ]")
	} else {
		color.New(color.FgHiRed).Print(username, " isn't a member of team [ "+teamName+" ]")
	}
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
		fmt.Println(i+1, "\t"+strconv.Itoa(team.ID), "\t"+team.Slug)
	}
}

func (mgr GitHubManager) MembershipOfTeams(username string) {

	color.New(color.Italic).Print(username + " Membership Of Teams\n")

	color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%23s\n", "No.", "ID", "TeamName")

	err, teams := mgr.Team.MembershipOfTeams(username)

	if err != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}
	for i, team := range teams {
		fmt.Printf("%3d\t%10d\t%23s\n", i+1, team.ID, team.Name)
	}
}

func (mgr GitHubManager) ListTeamMembers(option string) {
	start := time.Now()

	color.New(color.Italic).Print("Team members will include the members of child teams.\nTo list members in a team, the team must be visible to the authenticated user..\n")

	switch option {
	case "all":
		color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%23s\t%40s\t\tTeams\n", "No.", "ID", "Username", "Email")
	case "email":
		color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%23s\t%40s\n", "No.", "ID", "Username", "Email")
	case "team":
		color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%23s\t\tTeams\n", "No.", "ID", "Username")
	default:
		color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%23s\n", "No.", "ID", "Username")
	}

	//https://docs.github.com/en/rest/reference/teams#list-team-members
	err, i := mgr.Organization.ListOrgMember()
	if err != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}
	// load cache (GITHUB NOT SUPPROT API SO ,USE CACHE FOR IMPROVE PERFORMANCE)
	caches := mgr.loadCache()

	for i, member := range i {
		switch option {
		case "all":
			fmt.Printf("%3d\t%10d\t%23s\t%40s\t\t"+mgr.Team.MemberCacheByUser(caches, member.Login).Team+"\n", (i + 1), member.ID, member.Login, mgr.Team.MemberCacheByUser(caches, member.Login).Email)
		case "email":
			fmt.Printf("%3d\t%10d\t%23s\t%40s\n", i+1, member.ID, member.Login, mgr.Team.MemberCacheByUser(caches, member.Login).Email)
		case "team":
			fmt.Printf("%3d\t%10d\t%23s\t\t"+mgr.Team.MemberCacheByUser(caches, member.Login).Team+"\n", i+1, member.ID, member.Login)
		default:
			fmt.Printf("%3d\t%10d\t%23s\n", i+1, member.ID, member.Login)
		}
	}
	fmt.Println("\n----------------------------\nTime used is ", time.Since(start))
}

func (mgr GitHubManager) ShowListTeamMember(teamName string, role string, email string) {
	start := time.Now()

	//Header
	color.New(color.Italic).Print("Team members will include the members of child teams.\nTo list members in a [" + teamName + "] team, the team must be visible to the authenticated user.\n")

	if email == "show" {
		color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%23s\t%40s\n", "No.", "ID", "Username", "Email")
	} else {
		color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%23s\n", "No.", "ID", "Username")
	}

	// Process
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

	fmt.Println("\n----------------------------\nTime used is ", time.Since(start))
}

func (mgr GitHubManager) ShowListTeamMemberExclude(teamName string, teamExcude string, role string, email string) {
	start := time.Now()

	// Header
	color.New(color.Italic).Print("Team members will include the members of child teams.\nTo list members in a [" + teamName + "] team and Exclude [" + teamExcude + "] , the team must be visible to the authenticated user.\n")

	if email != "" {
		color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%23s\t%40s\n", "No.", "ID", "Username", "Email")
	} else {
		color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%23s\n", "No.", "ID", "Username")
	}

	// Process
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

	fmt.Println("\n----------------------------\nTime used is ", time.Since(start))
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

func (mgr GitHubManager) ReadProjectMemberListTemplateCSV(fileName string) {

	templ := csv.Template{}

	err, proj, csvTemplate := templ.ReadProjectMemberListTemplateCSV("reports/input/" + fileName)
	if err != nil {
		color.New(color.FgHiRed).Println(err.Error())
		os.Exit(1)
	}

	color.New(color.Italic).Print("CSV File Reader.\nTo list members in a  CSV file , [" + proj + "] team, the team must be visible to the GitHub.\n")

	color.New(color.FgHiMagenta).Printf("%2s\t%30s\t%40s\t%20s\t%10s\t%15s\n", "No.", "MemberName", "Email", "Role", "Team Role", "UserName")

	I := 0
	for _, csvTempl := range csvTemplate {
		if csvTempl.GitHub == "Y" {
			I++
			fmt.Printf("%2d\t%30s\t%40s\t%20s\t%10s\t%15s\n", I, csvTempl.Fullname, csvTempl.Email, csvTempl.Role, csvTempl.GitHubTeamRole, csvTempl.GitHubUsername)
		}
	}
}

func (mgr GitHubManager) InviteMemberToCorpTeamEmail(teamName string, email string) {
	color.New(color.Italic).Print("Create an organization invitation assign to [" + teamName + "] team. (org support member only) \n")

	//MEMBER ONLY!!
	const role string = "direct_member"

	// load cache (GITHUB NOT SUPPROT API ,SO WE USE CACHE FOR IMPROVE PERFORMANCE)
	caches := mgr.loadCache()

	// Invite member
	fmt.Printf(" %40s\t%20s : ", email, teamName)

	if mgr.User.CheckAlreadyMemberTeamByEmail(caches, email, teamName) {

		color.New(color.FgHiMagenta).Println("Already Exist")
		return
	}

	if mgr.IsInvited(teamName, email, "") {
		color.New(color.FgHiMagenta).Println("Invited")
		return
	}

	teamID := mgr.Team.GetInfoTeam(teamName).ID

	if err := mgr.Organization.InviteToCorpTeam(email, role, teamID); err != nil {
		color.New(color.FgHiRed).Println("Error ", err.Error())
		os.Exit(1)
	}

	caches = append(caches, model.Cache{Email: email})
	mgr.SetCache("cache/cache.csv", caches)

	color.New(color.FgHiGreen).Println("Done")
}

func (mgr GitHubManager) AddOrUpdateTeamMembershipUsername(teamName string, role string, username string) {
	color.New(color.Italic).Print("Add or update team membership for a user or Create an organization invitation assign to [" + teamName + "] team. \nAdds an organization member to a team., An authenticated organization owner or team maintainer can add organization members to a team. \n")
	mgr.AddOrUpdateTeamMembership(nil, "", teamName, role, username)
}

func (mgr GitHubManager) AddOrUpdateTeamMembershipEmail(teamName string, role string, email string) {
	color.New(color.Italic).Print("Add or update team membership for a user or Create an organization invitation assign to [" + teamName + "] team. \nAdds an organization member to a team., An authenticated organization owner or team maintainer can add organization members to a team. \n")

	// load cache (GITHUB NOT SUPPROT API ,SO WE USE CACHE FOR IMPROVE PERFORMANCE)
	caches := mgr.loadCache()

	mgr.AddOrUpdateTeamMembership(caches, email, teamName, role, "")
}

func (mgr GitHubManager) AddOrUpdateTeamMembership(caches []model.Cache, email string, teamName string, role string, username string) {

	if email != "" && username == "" {
		username = mgr.User.EmailToUsername(caches, email)
		if username == "" {
			color.New(color.FgYellow).Println("Email " + email + " Not Found in GitHub account or load Cache Please check again!")
			os.Exit(1)
		}
	}
	fmt.Printf(" %40s\t%20s\t%20s\t%20s : ", email, teamName, username, role)

	err, _ := mgr.AddOrUpdateTeamMembershipForAUser(username, teamName, role)
	if err != nil {
		color.New(color.FgHiRed).Println(err.Error())
		os.Exit(1)
	}

	color.New(color.FgHiGreen).Println("Done")
}

func (mgr GitHubManager) InviteMemberToCorpTeamTemplateCSVs() {
	color.New(color.Italic).Print("Create an organization invitation from input template files. \n")

	//get list fie name
	files, err := ioutil.ReadDir("reports/input/")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name(), " : ")
		mgr.inviteMemberToCorpTeamTemplateCSV(mgr.loadCache(), file.Name())
	}
}

func (mgr GitHubManager) InviteMemberToCorpTeamTemplateCSV(fileName string) {
	color.New(color.Italic).Print("Create an organization invitation from [" + fileName + "] file. \n")
	mgr.inviteMemberToCorpTeamTemplateCSV(mgr.loadCache(), fileName)
}

func (mgr GitHubManager) Debuging() {
	fmt.Println("Hello")

}

func (mgr GitHubManager) inviteMemberToCorpTeamTemplateCSV(caches []model.Cache, fileName string) {

	templ := csv.Template{}

	err, proj, csvTemplate := templ.ReadProjectMemberListTemplateCSV("reports/input/" + fileName)
	if err != nil {
		color.New(color.FgHiRed).Println(err.Error())
		os.Exit(1)
	}
	//Check Team in ORG
	if err, check := mgr.Team.CheckTeamInORG(proj); err != nil {
		color.New(color.FgHiRed).Println(err.Error())
		os.Exit(1)
	} else if !check {
		color.New(color.FgYellow).Println("Team [" + proj + "] Not Found in Organization !")
		os.Exit(1)
	}

	teamID := mgr.Team.GetInfoTeam(proj).ID

	I := 0

	for _, csvTempl := range csvTemplate {

		if csvTempl.GitHub == "Y" {
			if csvTempl.Email != "" && csvTempl.GitHubUsername != "" {
				I++
				fmt.Print(I, "\t")

				if mgr.IsInvited(proj, csvTempl.Email, csvTempl.GitHubUsername) {
					color.New(color.FgHiMagenta).Println("Invited")
					continue
				}

				mgr.AddOrUpdateTeamMembership(caches, csvTempl.Email, proj, csvTempl.GitHubTeamRole, csvTempl.GitHubUsername)
			} else if csvTempl.Email != "" {
				I++
				fmt.Print(I, "\t")

				// new Invite member via email
				fmt.Printf(" %40s\t%20s : ", csvTempl.Email, proj)

				if mgr.User.CheckAlreadyMemberTeamByEmail(caches, csvTempl.Email, proj) {
					color.New(color.FgHiMagenta).Println("Already Exist")
					continue
				}

				if mgr.IsInvited(proj, csvTempl.Email, "") {
					color.New(color.FgHiMagenta).Println("Invited")
					continue
				}

				//Invite API
				if err := mgr.Organization.InviteToCorpTeam(csvTempl.Email, "direct_member", teamID); err != nil {
					color.New(color.FgHiRed).Println("Error ", err.Error())
					os.Exit(1)
				}

				color.New(color.FgHiGreen).Println("Done")

				//mgr.inviteMemberToCorpTeam(caches, proj, "direct_member", csvTempl.Email)
			}
		} else if csvTempl.GitHub == "N" && csvTempl.GitHubUsername != "" {
			//reject out to team when N
			mgr.removeTeamMembershipForUser(caches, proj, csvTempl.GitHubUsername)

		}
	}
	mgr.SetCache("cache/cache.csv", caches)
}

func (mgr GitHubManager) ShowListTeamMemberPending(teamName string) {
	color.New(color.Italic).Print("List pending [" + teamName + "] team invitations\n")

	err, pendings := mgr.Organization.ListPendingTeamInvitations(teamName)
	if err != nil {
		color.New(color.FgHiRed).Println(err.Error())
		os.Exit(1)
	}

	color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%40s\t%20s\n", "No.", "ID", "Email", "Username")

	for i, invitation := range pendings {
		if invitation.Login != nil {
			fmt.Printf("%3d\t%10d\t%40s\t%20s\n", i+1, invitation.ID, invitation.Email, invitation.Login)
			continue
		}
		fmt.Printf("%3d\t%10d\t%40s\t%20s\n", i+1, invitation.ID, invitation.Email, "")
	}
}

func (mgr GitHubManager) ShowListPendingOrganizationInvitations() {
	color.New(color.Italic).Print("List pending organization invitations\n")

	err, pendings := mgr.Organization.ListPendingOrganizationInvitations()
	if err != nil {
		color.New(color.FgHiRed).Println(err.Error())
		os.Exit(1)
	}

	color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%40s\n", "No.", "ID", "Email")

	for i, invitation := range pendings {
		fmt.Printf("%3d\t%10d\t%40s\n", i+1, invitation.ID, invitation.Email)
	}
}

func (mgr GitHubManager) ExportCSVMemberTeams() {
	start := time.Now()
	if err, teams := mgr.Team.ListTeams(); err != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	} else {
		for _, team := range teams {
			mgr.ExportCSVMemberTeam(team.Slug)
		}
		color.New(color.FgHiGreen).Print("Done\n")
	}
	fmt.Println("\n----------------------------\nTime used is ", time.Since(start))
}

func (mgr GitHubManager) ExportCSVMemberTeam(teamName string) {

	color.New(color.Italic).Print("Export CSV Member Team [" + teamName + "] : ")

	var dataset []model.TeamMemberReport

	caches := mgr.loadCache()

	I := 0

	for _, role := range []string{"maintainer", "member"} {
		for _, teamMember := range mgr.Team.ListTeamMember(teamName, role) {
			I++
			dataset = append(dataset, model.TeamMemberReport{
				No:       strconv.Itoa(I),
				ID:       strconv.Itoa(teamMember.ID),
				Username: teamMember.Login,
				Name:     "",
				Email:    mgr.Team.MemberCacheByUser(caches, teamMember.Login).Email,
				Role:     role,
			})
		}
	}

	result := csv.WriteTeamMemberReport(teamName, "Report membership of team Generated by GHMGR "+mgr.Version+" : ", "reports/output/report-members-of-"+teamName, dataset)

	if result != nil {
		color.New(color.FgRed).Println(result.Error())
		os.Exit(1)
	}
	color.New(color.FgHiGreen).Print("Done\n")
}

func (mgr GitHubManager) ExportCSVMemberTeamTemplates() {

	err, teams := mgr.Team.ListTeams()
	if err != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}

	for i, team := range teams {

		if _, err := os.Stat("reports/input/" + team.Slug + ".csv"); !errors.Is(err, os.ErrNotExist) {
			fmt.Print(i+1, " : ")
			mgr.ExportCSVMemberTeamTemplate(team.Slug)
		}

	}
}

func (mgr GitHubManager) ExportCSVMemberTeamTemplate(teamName string) {

	color.New(color.Italic).Print("Export CSV Template Member Team [" + teamName + "] : ")
	var dataset []model.ProjectMemberListTemplate

	//Load Cache for get info
	caches := mgr.loadCache()

	//LOAD Templ
	templ := csv.Template{}
	err, _, csvTemplates := templ.ReadProjectMemberListTemplateCSV("reports/input/" + teamName + ".csv")
	if err != nil {
		color.New(color.FgHiRed).Println(err.Error())
		os.Exit(1)
	}

	var emails []string
	I := 0

	for _, role := range []string{"maintainer", "member"} {
		for _, teamMember := range mgr.Team.ListTeamMember(teamName, role) {
			I++
			email := mgr.Team.MemberCacheByUser(caches, teamMember.Login).Email

			templ := mgr.Team.CSVTemplate(csvTemplates, email)
			if templ.Email != "" {
				emails = append(emails, templ.Email)
			}

			dataset = append(dataset, model.ProjectMemberListTemplate{
				No:                strconv.Itoa(I),
				Username:          templ.Username,
				Fullname:          templ.Fullname,
				Email:             email,
				Role:              templ.Role,
				SubscriptionOwner: templ.SubscriptionOwner,
				GitHubUsername:    teamMember.Login,
				GitHubTeamRole:    role,
				GitHub:            "Y",
				AzureDEV:          templ.AzureDEV,
				AzurePRD:          templ.AzurePRD,
				ELK:               templ.ELK,
				Jumphost:          templ.Jumphost,
				Bastion:           templ.Bastion,
			})
		}
	}

	for _, csvTempl := range csvTemplates {
		if !mgr.User.CheckEmailInList(emails, csvTempl.Email) && csvTempl.Email != "" && csvTempl.Email != "verify email" {
			I++

			dataset = append(dataset, model.ProjectMemberListTemplate{
				No:                strconv.Itoa(I),
				Username:          csvTempl.Username,
				Fullname:          csvTempl.Fullname,
				Email:             csvTempl.Email,
				Role:              csvTempl.Role,
				SubscriptionOwner: csvTempl.SubscriptionOwner,
				GitHubUsername:    csvTempl.GitHubUsername,
				GitHubTeamRole:    csvTempl.GitHubTeamRole,
				GitHub:            csvTempl.GitHub,
				AzureDEV:          csvTempl.AzureDEV,
				AzurePRD:          csvTempl.AzurePRD,
				ELK:               csvTempl.ELK,
				Jumphost:          csvTempl.Jumphost,
				Bastion:           csvTempl.Bastion,
			})
		}
	}

	result := csv.Template{}.WriteProjectMemberListTemplateCSV(teamName, "template membership of team Generated by GHMGR "+mgr.Version+" : ", "reports/output/"+teamName, dataset)

	if result != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}
	color.New(color.FgHiGreen).Print("Done\n")
}

func (mgr GitHubManager) ExportCSVMemberTeamExclude(teamName string, teamExclude string) {
	start := time.Now()

	color.New(color.Italic).Print("Export CSV Member Team [" + teamName + "] Exclude Team [" + teamExclude + "] : ")
	var dataset []model.TeamMemberReport

	caches := mgr.loadCache()

	I := 0

	for _, role := range []string{"maintainer", "member"} {
		for _, teamMember := range mgr.Team.ListTeamMemberExcludeTeam(teamName, teamExclude, role) {
			I++
			dataset = append(dataset, model.TeamMemberReport{
				No:       strconv.Itoa(I),
				ID:       strconv.Itoa(teamMember.ID),
				Username: teamMember.Login,
				Name:     "",
				Email:    mgr.Team.MemberCacheByUser(caches, teamMember.Login).Email,
				Role:     role,
			})
		}
	}

	result := csv.WriteTeamMemberReport(teamName, "Report membership of team Generated by GHMGR "+mgr.Version+" : ", "reports/output/report-members-of-"+teamName, dataset)

	if result != nil {
		color.New(color.FgRed).Println(result.Error())
		os.Exit(1)
	}
	color.New(color.FgHiGreen).Print("Done\n")

	fmt.Println("\n----------------------------\nTime used is ", time.Since(start))

}

func (mgr GitHubManager) CancelOrganizationInvitationByEmail(email string) {
	color.New(color.Italic).Println("Cancel an organization invitation. In order to cancel an organization invitation, the authenticated user must be an organization owner.")

	caches := mgr.loadCache()

	mgr.cancelOrganizationInvitationByEmail(caches, email)

}

func (mgr GitHubManager) cancelOrganizationInvitationByEmail(caches []model.Cache, email string) {
	color.New(color.FgYellow).Print("Cancel invitation Email [" + email + "] : ")

	if err, invitationID := mgr.Organization.InviteEmailToInviteID(email); err != nil {
		color.New(color.FgHiRed).Println("ERROR ", err)
		os.Exit(1)
	} else {
		if err := mgr.Organization.CancelOrganizationInvitation(invitationID); err != nil {
			color.New(color.FgHiRed).Println("ERROR ", err)
			os.Exit(1)
		} else {
			var cs []model.Cache
			for _, c := range caches {
				if (c.Email != email) && c.Username != "" {
					cs = append(cs, c)
				}
			}
			mgr.SetCache("cache/cache.csv", cs)

			color.New(color.FgHiGreen).Println("Done")
		}
	}
}

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

func (mgr GitHubManager) ListDormantUsersfromCSV(filename string) {

	color.New(color.Italic).Println("List Dormant users of the organization from [" + filename + "] CSV file")

	color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%20s\t%20s\tTeams\n", "No.", "ID", "Username", "LastActive")

	err, dormantUsers := csv.Template{}.ReadDormantCSV("reports/input/" + filename)
	if err != nil {
		color.New(color.FgHiRed).Println(err.Error())
		os.Exit(1)
	}

	for i, dormantUser := range dormantUsers {

		ts := ""
		err, teams := mgr.Team.MembershipOfTeamsCacheTeam(dormantUser.Login)
		if err != nil {
			color.New(color.FgRed).Println(err.Error())
			os.Exit(1)
		}
		for j, team := range teams {
			if j == 0 {
				ts += team
			} else {
				ts += "|" + team
			}
		}
		fmt.Printf("%3d\t%10s\t%20s\t%20s\t"+ts+"\n", i+1, dormantUser.ID, dormantUser.Login, dormantUser.LastActive)
	}
}

func (mgr GitHubManager) ExportDormantUsersToCSV(filename string) {

	color.New(color.Italic).Println("Export report dormant users of the organization from [" + filename + "] CSV file")
	fmt.Print("Exporting : ")

	err, dormantUsers := csv.Template{}.ReadDormantCSV("reports/input/" + filename)
	if err != nil {
		color.New(color.FgHiRed).Println(err.Error())
		os.Exit(1)
	}

	var dataset []model.DormantUser
	for i, dormantUser := range dormantUsers {
		ts := ""
		err, teams := mgr.Team.MembershipOfTeamsCacheTeam(dormantUser.Login)
		if err != nil {
			color.New(color.FgRed).Println(err.Error())
			os.Exit(1)
		}
		for j, team := range teams {
			if j == 0 {
				ts += team
			} else {
				ts += "|" + team
			}
		}

		dataset = append(dataset, model.DormantUser{
			No:           strconv.Itoa(i + 1),
			CreateAt:     dormantUser.CreateAt,
			ID:           dormantUser.ID,
			Login:        dormantUser.Login,
			Role:         dormantUser.Role,
			Suspended:    dormantUser.Suspended,
			LastLoggedIP: dormantUser.LastLoggedIP,
			Dormant:      dormantUser.Dormant,
			LastActive:   dormantUser.LastActive,
			TwoFAEnabled: dormantUser.TwoFAEnabled,
			Teams:        ts,
			Excepted:     dormantUser.Excepted,
		})
	}

	result := csv.Template{}.WriteDormantCSV("reports/output/"+filename, dataset)
	if result != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}
	color.New(color.FgHiGreen).Println("Done")
}

func (mgr GitHubManager) RemoveDormantUsersFromCSV(backup bool, filename string) {

	color.New(color.Italic).Println("Remove dormant users of the organization from [" + filename + "] CSV file")

	if backup {
		color.New(color.FgCyan).Print("Backup Report : ")
		mgr.ExportDormantUsersToCSV(filename)
	}

	err, dormantUsers := csv.Template{}.ReadDormantCSV("reports/input/" + filename)
	if err != nil {
		color.New(color.FgHiRed).Println(err.Error())
		os.Exit(1)
	}

	I := 0
	for _, dormantUser := range dormantUsers {
		if dormantUser.Excepted != "Y" && dormantUser.Excepted != "Yes" {
			I++

			ts := ""
			err, teams := mgr.Team.MembershipOfTeamsCacheTeam(dormantUser.Login)
			if err != nil {
				color.New(color.FgRed).Println(err.Error())
				os.Exit(1)
			}
			for j, team := range teams {
				if j == 0 {
					ts += team
				} else {
					ts += "|" + team
				}
			}

			color.New(color.FgHiRed).Print(dormantUser.Login, " Removing : ")

			if err := mgr.Organization.RemoveOrganizationMember(dormantUser.Login); err != nil {
				color.New(color.FgHiRed).Println("ERROR ", err)
				os.Exit(1)
			}

			color.New(color.FgHiGreen).Println("Done")
		}

	}

}

//?
func (mgr GitHubManager) RemoveOrganizationMemberExculdeTeamMembers() {
	start := time.Now()

	color.New(color.Italic).Print("Remove a members will exclude the members out of child teams.\nTo list members out a team and remove membership of Org.\n")

	err, i := mgr.Organization.ListOrgMember()
	if err != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}

	caches := mgr.loadCache()

	I := 0
	for _, member := range i {
		if mgr.Team.CheckMembershipOutOfTeamsCache(caches, member.Login) {
			I++
			mgr.removeOrganizationMember(member.Login)
		}

	}
	fmt.Println("\n----------------------------\nTime used is ", time.Since(start))
}

func (mgr GitHubManager) RemoveOrganizationMember(username string) {

	color.New(color.Italic).Print("Removing a user from this list will remove them from all teams and they will no longer have any access to the organization's repositories\n")

	mgr.removeOrganizationMember(username)
}

// USE PLEASE CONERT !
func (mgr GitHubManager) RemoveOrganizationMembersWithoutTeam() {
	start := time.Now()
	color.New(color.Italic).Print("Removing a users from Oranization in Condition with Team is Null\n")

	i := 0

	for _, c := range mgr.loadCache() {
		if c.Team == "" {
			i++
			fmt.Print(i, " : ")
			mgr.removeOrganizationMember(c.Username)
		}
	}
	fmt.Println("\n----------------------------\nTime used is ", time.Since(start))
}

// USE PLEASE CONCERT !
func (mgr GitHubManager) RemoveOrganizationMembersWithoutMembershipOfTeams() {
	start := time.Now()
	color.New(color.Italic).Print("Removing a users from Oranization without membership of teams\n")

	i := 0
	//check team null
	for _, c := range mgr.loadCache() {
		if c.Team == "" {
			i++
			fmt.Print(i, " : ")

			mgr.removeOrganizationMember(c.Username)
		}
	}

	//BACKUP REMOVE USER..

	fmt.Println("\n----------------------------\nTime used is ", time.Since(start))
}

func (mgr GitHubManager) RemoveMembershipOfTeamWithoutEmail(team string) {
	start := time.Now()
	color.New(color.Italic).Print("Remove a membership of team don't verify email\n")

	i := 0
	cache := mgr.loadCache()
	for _, c := range cache {
		if c.Email == "" && isMembershipOfTeamInCache(c.Team, team) {
			i++
			fmt.Print(i, " : ")
			mgr.removeTeamMembershipForUser(cache, team, c.Username)
		}
	}
	fmt.Println("\n----------------------------\nTime used is ", time.Since(start))
}

func isMembershipOfTeamInCache(teamsString string, team string) bool {
	for _, t := range strings.Split(teamsString, "|") {
		if t == team {
			return true
		}
	}
	return false
}

// !
func (mgr GitHubManager) RemoveOrganizationMembersWithoutEmail() {
	start := time.Now()
	color.New(color.Italic).Print("Removing a users from Oranization without Email\n")

	i := 0
	//check email empty
	for _, c := range mgr.loadCache() {
		if c.Email == "" {
			i++
			fmt.Print(i, " : ")
			// debug only
			color.New(color.FgHiRed).Print(c.Username, " xx removing an organization : ")
			color.New(color.FgHiMagenta).Println(" Done")
			//mgr.removeOrganizationMember(c.Username)

		}
	}
	fmt.Println("\n----------------------------\nTime used is ", time.Since(start))
}

func (mgr GitHubManager) RemoveOrganizationMembers() {
	start := time.Now()
	color.New(color.Italic).Print("Removing a users from Oranization in Condition with Email is Empty or Team is Null\n")

	i := 0
	//check email empty and team null
	for _, c := range mgr.loadCache() {
		if c.Email == "" || c.Team == "" {
			i++
			fmt.Print(i, " : ")
			//mgr.removeOrganizationMember(c.Username)
		}
	}
	fmt.Println("\n----------------------------\nTime used is ", time.Since(start))
}

func (mgr GitHubManager) removeOrganizationMember(username string) {
	color.New(color.FgHiRed).Print(username, " removing an organization : ")

	if err := mgr.Organization.RemoveOrganizationMember(username); err != nil {
		color.New(color.FgHiRed).Println("ERROR ", err)
		os.Exit(1)
	}

	color.New(color.FgHiGreen).Println(" Done")
}

func (mgr GitHubManager) RemoveTeamMembershipForUser(teamname, username string) {
	color.New(color.Italic).Print("To remove a membership between a user and a team, the authenticated user must have 'admin' permissions to the team or be an owner of the organization that the team is associated with. Removing team membership does not delete the user, it just removes their membership from the team.\n")

	mgr.removeTeamMembershipForUser(mgr.loadCache(), teamname, username)
}

func (mgr GitHubManager) removeTeamMembershipForUser(cache []model.Cache, teamname, username string) {
	color.New(color.FgHiRed).Print(username, " removing a "+teamname+" team :")
	if err := mgr.Team.RemoveTeamMembershipForUser(teamname, username); err != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}
	color.New(color.FgHiGreen).Println(" Done")

	mgr.RemoveMemberCachedInvited(teamname, mgr.UsernameToEmail(cache, username))
}

func (mgr GitHubManager) ListExculdeTeamMembers() {
	start := time.Now()

	color.New(color.Italic).Print("Team members will exclude the members of child teams.\nTo list members out a team, the team must be visible to the authenticated user..\n")

	color.New(color.FgHiMagenta).Printf("%3s\t%10s\t%23s\t%40s\n", "No.", "ID", "Username", "Email")

	//https://docs.github.com/en/rest/reference/teams#list-team-members
	err, i := mgr.Organization.ListOrgMember()
	if err != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}
	// load cache (GITHUB NOT SUPPROT API SO ,USE CACHE FOR IMPROVE PERFORMANCE)
	caches := mgr.loadCache()

	I := 0
	for _, member := range i {
		if mgr.Team.CheckMembershipOutOfTeamsCache(caches, member.Login) {
			I++
			fmt.Printf("%3d\t%10d\t%23s\t%40s\n", I, member.ID, member.Login, mgr.Team.MemberCacheByUser(caches, member.Login).Email)
		}

	}
	fmt.Println("\n----------------------------\nTime used is ", time.Since(start))
}

func (mgr GitHubManager) ExportORGMemberReport() {
	start := time.Now()
	color.New(color.Italic).Print("Export CSV Members of organization : ")

	os.Mkdir("reports", 0755)
	os.Mkdir("reports/output", 0755)

	//https://docs.github.com/en/rest/reference/teams#list-team-members
	err, i := mgr.Organization.ListOrgMember()
	if err != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}
	// load cache (GITHUB NOT SUPPROT API SO ,USE CACHE FOR IMPROVE PERFORMANCE)
	caches := mgr.loadCache()

	var dataset []model.OrgMemberReport
	for i, member := range i {
		ds := model.OrgMemberReport{
			No:       strconv.Itoa(i + 1),
			ID:       strconv.Itoa(member.ID),
			Username: member.Login,
			Name:     "",
			Email:    mgr.Team.MemberCacheByUser(caches, member.Login).Email,
			Team:     mgr.Team.MemberCacheByUser(caches, member.Login).Team,
		}

		dataset = append(dataset, ds)
	}
	result := csv.WriteORGMemberReport("Report membership of organization Generated by GHMGR "+mgr.Version+" : ", "reports/output/report-members-of-organization", dataset)

	if result != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}
	color.New(color.FgHiGreen).Print("Done\n")

	fmt.Println("\n----------------------------\nTime used is ", time.Since(start))
}

func (mgr GitHubManager) ExportORGMemberWithOutMembershipOfTeamReport() {
	start := time.Now()
	color.New(color.Italic).Print("Export CSV members of organization without membership of team(s) : ")

	os.Mkdir("reports", 0755)
	os.Mkdir("reports/output", 0755)

	//https://docs.github.com/en/rest/reference/teams#list-team-members
	err, i := mgr.Organization.ListOrgMember()
	if err != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}
	// load cache (GITHUB NOT SUPPROT API SO ,USE CACHE FOR IMPROVE PERFORMANCE)
	caches := mgr.loadCache()

	I := 0
	var dataset []model.OrgMemberReport
	for i, member := range i {
		if mgr.Team.CheckMembershipOutOfTeamsCache(caches, member.Login) {
			I++
			ds := model.OrgMemberReport{
				No:       strconv.Itoa(i + 1),
				ID:       strconv.Itoa(member.ID),
				Username: member.Login,
				Name:     "",
				Email:    mgr.Team.MemberCacheByUser(caches, member.Login).Email,
				Team:     mgr.Team.MemberCacheByUser(caches, member.Login).Team,
			}

			dataset = append(dataset, ds)
		}

	}
	result := csv.WriteORGMemberReport("Report membership of organization without membership of team(s) Generated by GHMGR "+mgr.Version+" : ", "reports/output/report-members-of-organization-without-membership-of-teams", dataset)

	if result != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}
	color.New(color.FgHiGreen).Print("Done\n")

	fmt.Println("\n----------------------------\nTime used is ", time.Since(start))
}
