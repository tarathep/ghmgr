package manage

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/tarathep/ghmgr/csv"
	"github.com/tarathep/ghmgr/model"
)

func (mgr GitHubManager) CachePending(teamName string) {

	err, pendings := mgr.Organization.ListPendingTeamInvitations(teamName)
	if err != nil {
		color.New(color.FgHiRed).Println(err.Error())
		os.Exit(1)
	}

	cache_pendings := []model.CachePending{}

	for i, invitation := range pendings {

		if invitation.Login != nil {
			cache_pendings = append(cache_pendings, model.CachePending{
				No:       strconv.Itoa(i + 1),
				ID:       strconv.Itoa(invitation.ID),
				Email:    invitation.Email,
				Username: fmt.Sprintf("%v", invitation.Login),
			})
			continue
		}
		cache_pendings = append(cache_pendings, model.CachePending{
			No:       strconv.Itoa(i + 1),
			ID:       strconv.Itoa(invitation.ID),
			Email:    invitation.Email,
			Username: "",
		})
	}

	// for _, cache_pending := range cache_pendings {
	// 	fmt.Println(cache_pending)
	// }
	csv.Pending{}.WriteCache("./cache/invited/"+teamName+".csv", cache_pendings)
}

func (mgr GitHubManager) IsInvited(teamName string, email string, username string) bool {
	csvFileName := "./cache/invited/" + teamName + ".csv"

	//check file exist?
	if _, err := os.Stat(csvFileName); errors.Is(err, os.ErrNotExist) {
		// does not exist file csv team
		// init file
		mgr.CachePending(teamName)
	}

	//Read members
	err, cache_pendings := csv.Pending{}.ReadCache("./cache/invited/" + teamName + ".csv")
	if err != nil {
		color.New(color.FgHiRed).Println(err.Error())
		os.Exit(1)
	}

	if isMemberInPendingCached(email, cache_pendings) {
		return true
	} else {
		cache_pendings = append(cache_pendings, model.CachePending{
			No:       strconv.Itoa(len(cache_pendings) + 1),
			ID:       "",
			Email:    email,
			Username: username,
		})
	}
	csv.Pending{}.WriteCache("./cache/invited/"+teamName+".csv", cache_pendings)
	return false
}

func isMemberInPendingCached(email string, cache_pendings []model.CachePending) bool {
	//Check member is exist
	for _, cache_pending := range cache_pendings {
		if cache_pending.Email == email && email != "" {
			return true
		}
	}
	return false
}

func (mgr GitHubManager) RemoveMemberCachedInviteds(teamName, emails string) {
	for _, email := range strings.Split(emails, ",") {
		mgr.RemoveMemberCachedInvited(teamName, email)
	}
}

func (mgr GitHubManager) RemoveMemberCachedInvited(teamName string, email string) {
	csvFileName := "./cache/invited/" + teamName + ".csv"

	if _, err := os.Stat(csvFileName); errors.Is(err, os.ErrNotExist) {
		mgr.CachePending(teamName)
	}

	new_cache_pendings := []model.CachePending{}

	err, cache_pendings := csv.Pending{}.ReadCache("./cache/invited/" + teamName + ".csv")
	if err != nil {
		color.New(color.FgHiRed).Println(err.Error())
		os.Exit(1)
	}

	for _, cache_pending := range cache_pendings {
		if cache_pending.Email != email {
			new_cache_pendings = append(new_cache_pendings, cache_pending)
		}
	}

	csv.Pending{}.WriteCache("./cache/invited/"+teamName+".csv", new_cache_pendings)
}
