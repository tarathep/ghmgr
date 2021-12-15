package main

import (
	"github.com/tarathep/githuby/github"
	"github.com/tarathep/githuby/login"
)

func main() {
	auth := login.Auth{}
	auth.Token = auth.GetToken()

	team := github.Team{Auth: auth, Debug: true}

	//team.GetRepoList("corp-ais", "myChannel")

	//team.UpdateRepoPermissionTeam("admin", "myChannel", "corp-ais", "demo-pipeline")
	team.AddnewTeamInAnotherRepoTeam("corp-ais", "IBM", "myChannel", "admin")
	// role string, teamName string, owner string, repoName string
}
