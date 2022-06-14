package manage

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/tarathep/ghmgr/model"
)

// load cache (GITHUB NOT SUPPROT API ,SO WE USE CACHE FOR IMPROVE PERFORMANCE)
func (mgr GitHubManager) loadCache() []model.Cache {
	err, caches := mgr.GetCache("./cache/cache.csv")
	if err != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}
	return caches
}

func (mgr GitHubManager) Caching() {
	color.New(color.Italic).Println("Cache GitHub members of the organization.")
	start := time.Now()
	os.Mkdir("cache", 0755)
	os.Mkdir("cache/teams", 0755)

	mgr.ExportTeamMemberCache()
	mgr.ExportMembersORGCache()
	fmt.Println("\n----------------------------\nTime used is ", time.Since(start))
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func (mgr GitHubManager) ExportTeamMemberCache() {

	removeContents("cache/teams")
	err, teams := mgr.Team.ListTeams()
	if err != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}

	color.New(color.FgHiCyan).Print("Caching Membership of Teams : ")
	for _, team := range teams {
		var cache []model.Cache
		for i, teamMember := range mgr.Team.ListTeamMember(team.Slug, "") {
			cache = append(cache, model.Cache{
				No:       strconv.Itoa(i + 1),
				ID:       strconv.Itoa(teamMember.ID),
				Username: teamMember.Login,
				Email:    "",
				Team:     team.Slug,
			})
		}
		mgr.SetCache("cache/teams/"+team.Slug+".csv", cache)
	}
	color.New(color.FgHiGreen).Print("Done\n")
}

func (mgr GitHubManager) ExportMembersORGCache() {
	var cache []model.Cache

	color.New(color.FgHiCyan).Print("Caching Member of Organization : ")

	err, listOrgMember := mgr.Organization.ListOrgMember()
	if err != nil {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}

	for i, orgMem := range listOrgMember {

		//GET MEMBER
		var membership string
		err, teams := mgr.Team.MembershipOfTeamsCacheTeam(orgMem.Login)
		if err != nil {
			fmt.Println(err)
		}
		for i, team := range teams {
			if i == len(teams)-1 {
				membership += team
			} else {
				membership += team + "|"
			}

		}

		_, usr := mgr.UserInfo(orgMem.Login)

		cache = append(cache, model.Cache{
			No:       strconv.Itoa(i + 1),
			ID:       strconv.Itoa(orgMem.ID),
			Username: orgMem.Login,
			Email:    usr.Email,
			Team:     membership,
		})
	}

	//save
	mgr.SetCache("cache/cache.csv", cache)
	color.New(color.FgHiGreen).Print("Done\n")
}
