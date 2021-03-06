package model

import "time"

type Invitation struct {
	ID           int         `json:"id"`
	NodeID       string      `json:"node_id"`
	Login        interface{} `json:"login"`
	Email        string      `json:"email"`
	Role         string      `json:"role"`
	CreatedAt    time.Time   `json:"created_at"`
	FailedAt     interface{} `json:"failed_at"`
	FailedReason interface{} `json:"failed_reason"`
	Inviter      struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"inviter"`
	TeamCount          int    `json:"team_count"`
	InvitationTeamsURL string `json:"invitation_teams_url"`
}
