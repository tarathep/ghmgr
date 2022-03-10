package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/tarathep/ghmgr/github"
	"github.com/tarathep/ghmgr/login"
	"github.com/tarathep/ghmgr/manage"
)

type Options struct {
	Name    string `short:"n" long:"name" description:"a name"`
	Email   string `short:"m" long:"email" description:"Email"`
	File    string `short:"f" long:"file" description:"File"`
	Team    string `short:"t" long:"team" description:"Team"`
	Pending bool   `short:"p" long:"pending" description:"Pending"`
	Role    string `short:"r" long:"role" description:"Role"`
	Version bool   `short:"v" long:"version" description:"Version"`
	Token   string `long:"token" description:"Personal Access Token"`
}

func main() {

	// teamID := team.GetInfoTeam("ipfm").ID
	// fmt.Print(teamID)

	// team.GetRepoList("myChannel")

	//team.UpdateRepoPermissionTeam("admin", "myChannel", "corp-ais", "demo-pipeline")
	// team.AddnewTeamInAnotherRepoTeam("corp-ais", "IBM", "myChannel", "admin")
	// role string, teamName string, owner string, repoName string

	// member.InviteToCorpTeam("bokie.demo@gmail.com", "direct_member", teamID)
	// member.ListPendingTeamInvitations("ipfm")

	// for i, teamMember := range team.ListTeamMember("enterprise-solution-development") {
	// 	println(i, teamMember.Login)
	// }

	//--- Options Flags ---
	var options Options
	parser := flags.NewParser(&options, flags.PrintErrors|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		log.Fatal(err)
	}

	flags.NewIniParser(parser)

	if options.Version {
		fmt.Print("v1.0.0")
	}

	if len(os.Args) > 1 {
		// SET AUTH & CORP
		auth := login.Auth{}
		auth.Token = auth.GetToken()

		team := github.Team{Auth: auth, Owner: "corp-ais"}
		member := github.Member{Auth: auth, Owner: "corp-ais"}
		gitHubMgr := manage.GitHubManager{Member: member, Team: team}

		switch os.Args[1] {

		case "list":
			{
				if len(os.Args) > 2 && os.Args[2] == "member" {

					if options.Team != "" && options.Pending {
						gitHubMgr.ShowListTeamMemberPending(options.Team)
					} else if options.Team != "" && options.Role != "" {
						gitHubMgr.ShowListTeamMember(options.Team, options.Role)
					} else if options.Team != "" {
						gitHubMgr.ShowListTeamMember(options.Team, "all")
					}

					if options.File != "" {
						gitHubMgr.ReadCSVFile(options.File)
					}
				}
			}
		case "export":
			{
				if len(os.Args) > 2 && os.Args[2] == "member" {
					if options.Team != "" {
						gitHubMgr.ExportCSVMemberTeam(options.Team)
					}
				}
			}
		case "invite":
			{
				if len(os.Args) > 2 && os.Args[2] == "member" {
					if options.Team != "" && options.Email != "" && options.Role != "" {
						// fmt.Print("InviteMemberToCorpTeam")
						gitHubMgr.InviteMemberToCorpTeam(options.Team, options.Role, options.Email)
					}
					if options.File != "" {
						// fmt.Print("InviteMemberToCorpTeamCSV")
						gitHubMgr.InviteMemberToCorpTeamCSV(options.File)
					}
				}
			}
		case "login":
			{
				if options.Token != "" {
					auth.SetToken(options.Token)
				}
			}

		}
	}
}
