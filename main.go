package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/tarathep/githuby/login"
	"github.com/tarathep/githuby/model"
)

func main() {
	auth := login.Auth{}
	auth.Token = auth.GetToken()

	// team := github.Team{Auth: auth, Owner: "corp-ais"}
	// teamID := team.GetInfoTeam("ccsm").ID
	// fmt.Print(teamID)

	// team.GetRepoList("myChannel")

	//team.UpdateRepoPermissionTeam("admin", "myChannel", "corp-ais", "demo-pipeline")
	// team.AddnewTeamInAnotherRepoTeam("corp-ais", "IBM", "myChannel", "admin")
	// role string, teamName string, owner string, repoName string

	// member := github.Member{Auth: auth, Debug: true, Owner: "corp-ais"}
	// member.InviteToCorpTeam("bokie.demo@gmail.com", "direct_member", teamID)

	records, err := readData("mc.csv")

	if err != nil {
		log.Fatal(err)
	}

	for i, record := range records {
		switch i {
		case 0:
			fmt.Println("TeamName/Project :>", record[1])
		case 1:
		default:
			csv := model.CSV{record[0], record[1], record[2], record[3], record[4], record[5]}
			fmt.Println((i - 1), csv.Email)
		}

	}
}

func readData(fileName string) ([][]string, error) {

	f, err := os.Open(fileName)

	if err != nil {
		return [][]string{}, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}
