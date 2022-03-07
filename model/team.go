package model

import "time"

type Team struct {
	Name            string    `json:"name"`
	ID              int       `json:"id"`
	NodeID          string    `json:"node_id"`
	Slug            string    `json:"slug"`
	Description     string    `json:"description"`
	Privacy         string    `json:"privacy"`
	URL             string    `json:"url"`
	HTMLURL         string    `json:"html_url"`
	MembersURL      string    `json:"members_url"`
	RepositoriesURL string    `json:"repositories_url"`
	Permission      string    `json:"permission"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	MembersCount    int       `json:"members_count"`
	ReposCount      int       `json:"repos_count"`
	Organization    struct {
		Login                   string      `json:"login"`
		ID                      int         `json:"id"`
		NodeID                  string      `json:"node_id"`
		URL                     string      `json:"url"`
		ReposURL                string      `json:"repos_url"`
		EventsURL               string      `json:"events_url"`
		HooksURL                string      `json:"hooks_url"`
		IssuesURL               string      `json:"issues_url"`
		MembersURL              string      `json:"members_url"`
		PublicMembersURL        string      `json:"public_members_url"`
		AvatarURL               string      `json:"avatar_url"`
		Description             string      `json:"description"`
		Name                    string      `json:"name"`
		Company                 interface{} `json:"company"`
		Blog                    string      `json:"blog"`
		Location                string      `json:"location"`
		Email                   interface{} `json:"email"`
		TwitterUsername         interface{} `json:"twitter_username"`
		IsVerified              bool        `json:"is_verified"`
		HasOrganizationProjects bool        `json:"has_organization_projects"`
		HasRepositoryProjects   bool        `json:"has_repository_projects"`
		PublicRepos             int         `json:"public_repos"`
		PublicGists             int         `json:"public_gists"`
		Followers               int         `json:"followers"`
		Following               int         `json:"following"`
		HTMLURL                 string      `json:"html_url"`
		CreatedAt               time.Time   `json:"created_at"`
		UpdatedAt               time.Time   `json:"updated_at"`
		Type                    string      `json:"type"`
	} `json:"organization"`
	Parent interface{} `json:"parent"`
}
