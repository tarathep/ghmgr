package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/tarathep/githuby/login"
)

type Team struct {
	Auth login.Auth
}

func (team Team) GetRepos() {
	req, err := http.NewRequest("GET", "https://api.github.com/orgs/corp-ais/teams/example/repos", nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	//req.Header.Set("Authorization", "token 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")
	req.Header.Set("Authorization", "token "+team.Auth.Token())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//bodyString := string(bodyBytes)

	//fmt.Println(bodyString)

	repos := Repos{}

	json.Unmarshal(bodyBytes, &repos)

	fmt.Println(repos[0].Name)
}
