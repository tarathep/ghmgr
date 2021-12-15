package main

import (
	"github.com/tarathep/githuby/github"
	"github.com/tarathep/githuby/login"
)

func main() {
	auth := login.Auth{}

	team := github.Team{Auth: auth}

	team.GetRepos()

}
