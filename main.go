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
	Name     string `short:"n" long:"name" description:"a name"`
	Email    string `short:"m" long:"email" description:"Email"`
	Username string `short:"u" long:"username" description:"Username GitHub"`
	File     string `short:"f" long:"file" description:"File"`
	Team     string `short:"t" long:"team" description:"Team"`
	Pending  bool   `short:"p" long:"pending" description:"Pending"`
	Cancel   bool   `short:"c" long:"cancel" description:"Cancel"`
	Role     string `short:"r" long:"role" description:"Role"`
	Version  bool   `short:"v" long:"version" description:"Version"`
	Token    string `long:"token" description:"Personal Access Token"`
	Owner    string `long:"owner" description:"Owner of GitHub"`
	ORG      bool   `short:"o" long:"org" description:"ORG"`
	Exclude  string `short:"e" long:"exclude" description:"Excude Team"`
	ID       string `long:"id" description:"ID for references such as GitHub ID, Personal etc."`
	Dormant  string `short:"d" long:"dormant" description:"Dormant Users in ORG & Teams"`
	Backup   bool   `short:"b" long:"backup" description:"Backup file or Report"`
}

const version string = "v1.3.1"

func main() {

	// Init Directory
	os.Mkdir("reports", 0755)
	os.Mkdir("reports/input", 0755)
	os.Mkdir("reports/output", 0755)

	//--- Options Flags ---
	var options Options
	parser := flags.NewParser(&options, flags.PrintErrors|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		log.Fatal(err)
	}

	flags.NewIniParser(parser)

	if options.Version {
		fmt.Print(version)
	}

	if len(os.Args) > 1 {
		// SET AUTH & CORP
		auth := login.Auth{}

		auth.Token = auth.GetToken()
		if auth.Token == "" {
			auth.Token = options.Token
		}
		auth.Owner = auth.GetOwner()
		if auth.Owner == "" {
			auth.Owner = options.Owner
		}

		team := github.Team{Auth: auth, Owner: auth.Owner}
		organization := github.Organization{Auth: auth, Owner: auth.Owner}
		user := github.User{Auth: auth}

		gitHubMgr := manage.GitHubManager{Version: version, Organization: organization, Team: team, User: user}

		switch os.Args[1] {

		case "list":
			{
				if len(os.Args) > 2 && os.Args[2] == "team" {
					if options.Username != "" {
						gitHubMgr.MembershipOfTeams(options.Username)
					} else {
						gitHubMgr.ListTeam()
					}

				} else if len(os.Args) > 2 && os.Args[2] == "member" {

					if len(os.Args) > 3 && os.Args[3] == "dormant" {
						if options.File != "" {
							gitHubMgr.ListDormantUsersfromCSV(options.File)
							return
						}
					}

					if options.Team != "" && options.Exclude != "" {
						if options.Pending {
							gitHubMgr.ShowListTeamMemberPending(options.Team)
						} else if options.Role != "" {
							gitHubMgr.ShowListTeamMemberExclude(options.Team, options.Exclude, options.Role, options.Email)
						} else if options.Email == "show" {
							gitHubMgr.ShowListTeamMemberExclude(options.Team, options.Exclude, "all", options.Email)
						} else {
							gitHubMgr.ShowListTeamMemberExclude(options.Team, options.Exclude, "all", options.Email)
						}
					} else if options.Team != "" && options.Team != "show" {
						if options.Pending {
							gitHubMgr.ShowListTeamMemberPending(options.Team)
						} else if options.Role != "" {
							gitHubMgr.ShowListTeamMember(options.Team, options.Role, options.Email)
						} else if options.Email == "show" {
							gitHubMgr.ShowListTeamMember(options.Team, "all", options.Email)
						} else {
							gitHubMgr.ShowListTeamMember(options.Team, "all", options.Email)
						}
					}

					if options.File != "" {
						// gitHubMgr.ReadCSVFile(options.File)
						gitHubMgr.ReadProjectMemberListTemplateCSV(options.File)
					}

					if options.ORG && options.Email == "show" && options.Team == "show" {
						gitHubMgr.ListTeamMembers("all")
					} else if options.ORG && options.Exclude == "team" {
						gitHubMgr.ListExculdeTeamMembers()
					} else if options.ORG && options.Email == "show" {
						gitHubMgr.ListTeamMembers("email")
					} else if options.ORG && options.Team == "show" {
						gitHubMgr.ListTeamMembers("team")
					} else if options.ORG && options.Pending {
						gitHubMgr.ShowListPendingOrganizationInvitations()
					} else if options.ORG {
						gitHubMgr.ListTeamMembers("")
					}

				}
			}
		case "export":
			{
				if len(os.Args) > 2 && os.Args[2] == "template" {
					if options.Team == "all" {
						gitHubMgr.ExportCSVMemberTeamTemplates()
						return
					}

					if options.Team != "" {
						gitHubMgr.ExportCSVMemberTeamTemplate(options.Team)
						return
					}

				}
				if len(os.Args) > 2 && os.Args[2] == "member" {

					if len(os.Args) > 3 && os.Args[3] == "dormant" {
						if options.File != "" {
							gitHubMgr.ExportDormantUsersToCSV(options.File)
							return
						}
					}

					if options.Team != "" && options.Exclude != "" {
						gitHubMgr.ExportCSVMemberTeamExclude(options.Team, options.Exclude)
					} else if options.Team == "all" {
						gitHubMgr.ExportCSVMemberTeams()
					} else if options.Team != "" {
						gitHubMgr.ExportCSVMemberTeam(options.Team)
					} else if options.ORG && options.Exclude == "team" {
						gitHubMgr.ExportORGMemberWithOutMembershipOfTeamReport()
					} else if options.ORG {
						gitHubMgr.ExportORGMemberReport()
					}

				}
			}
		case "import":
			{
				if len(os.Args) > 2 && os.Args[2] == "template" {
					if options.File == "all" {
						gitHubMgr.InviteMemberToCorpTeamTemplateCSVs()
						return
					}
					if options.File != "" {
						gitHubMgr.InviteMemberToCorpTeamTemplateCSV(options.File)
						return
					}
				}
			}
		case "invite":
			{
				if len(os.Args) > 2 && os.Args[2] == "member" {
					if options.Team != "" && options.Email != "" {
						gitHubMgr.InviteMemberToCorpTeamEmail(options.Team, options.Email)
						return
					}
					if options.Cancel && options.ID != "" {
						gitHubMgr.CancelOrganizationInvitation(options.ID)
						return
					}
					if options.Cancel && options.Email != "" {
						gitHubMgr.CancelOrganizationInvitationByEmail(options.Email)
						return
					}
				}
			}
		case "create":
			{
				if len(os.Args) > 2 && os.Args[2] == "team" {
					//==>
				}
			}
		case "add":
			{
				if len(os.Args) > 2 && os.Args[2] == "member" {
					if options.Team != "" && options.Username != "" && options.Role != "" {
						gitHubMgr.AddOrUpdateTeamMembershipUsername(options.Team, options.Role, options.Username)
						return
					}
					if options.Team != "" && options.Email != "" && options.Role != "" {
						gitHubMgr.AddOrUpdateTeamMembershipEmail(options.Team, options.Role, options.Email)
						return
					}
				}
			}
		case "remove":
			{
				if len(os.Args) > 2 && os.Args[2] == "member" {

					if options.Username != "" && options.ORG {
						gitHubMgr.RemoveOrganizationMember(options.Username)
						return
					}
					if options.ORG && options.Exclude == "team" {
						gitHubMgr.RemoveOrganizationMemberExculdeTeamMembers()
						return
					}

					if options.ORG && options.Team == "null" {
						gitHubMgr.RemoveOrganizationMembersWithoutMembershipOfTeams()
						return
					}

					if options.ORG && options.Email == "null" {
						gitHubMgr.RemoveOrganizationMembersWithoutEmail()
						return
					}

					if options.ORG {
						//remove all with condition email = null and team null
						//gitHubMgr.RemoveOrganizationMembers()
						return
					}

					if options.Username != "" && options.Team != "" {
						gitHubMgr.RemoveTeamMembershipForUser(options.Team, options.Username)
						return
					}
					if len(os.Args) > 3 && os.Args[3] == "dormant" {
						if options.File != "" {
							gitHubMgr.RemoveDormantUsersFromCSV(options.Backup, options.File)
							return
						}
					}
				}
			}
		case "login":
			{
				if options.Token != "" {
					auth.SetToken(options.Token)
				}
				if options.Owner != "" {
					auth.SetOwner(options.Owner)
				}

			}
		case "load":
			{
				if len(os.Args) > 2 && os.Args[2] == "cache" {
					gitHubMgr.Caching()
				}
			}
		case "check":
			{
				if len(os.Args) > 2 && os.Args[2] == "member" {
					if options.Username != "" && options.Team != "" {
						gitHubMgr.CheckTeamMembershipForUser(options.Team, options.Username)
					} else if options.Username != "" {
						gitHubMgr.CheckOrganizationMembership(options.Username)
					}
				}
			}
		case "get":
			{
				if len(os.Args) > 2 && os.Args[2] == "member" {
					if len(os.Args) > 3 && os.Args[3] == "username" {
						if options.Email != "" {
							gitHubMgr.GetUsernameFromEmail(options.Email)
							return
						}
					}
					if len(os.Args) > 3 && os.Args[3] == "email" {
						if options.Username != "" {
							gitHubMgr.GetEmailFromUsername(options.Username)
							return
						}
					}
				}
			}
		}
	}
}
