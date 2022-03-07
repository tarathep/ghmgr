package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/tarathep/githuby/login"
)

type Member struct {
	Auth  login.Auth
	Debug bool
	Owner string
}

func (member Member) InviteToCorpTeam(Email string, Role string, teamID int) {

	type Payload struct {
		// Invitee_ID int    `json:"invitee_id"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		Team_IDs []int  `json:"team_ids"`
	}

	data := Payload{
		Email:    Email,
		Role:     Role,
		Team_IDs: []int{teamID},
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api.github.com/orgs/"+member.Owner+"/invitations", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "token "+member.Auth.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		bodyString := string(bodyBytes)
		log.Fatal(resp.Status, bodyString)
	}

	bodyString := string(bodyBytes)
	fmt.Print(bodyString)

	// teamm := model.Team{}
	// json.Unmarshal(bodyBytes, &teamm)
}
