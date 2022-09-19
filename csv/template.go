package csv

import (
	"encoding/csv"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/tarathep/ghmgr/model"
)

type Template struct{}

func (Template) ReadFile(filename string) (err error, project string, csvlist []model.CSV) {

	records, err := readData(filename)

	if err != nil {
		log.Fatal(err)
		return err, "", nil
	}

	for i, record := range records {
		switch i {
		case 0:
			reg, err := regexp.Compile("[^A-Za-z0-9]+")
			if err != nil {
				log.Fatal(err)
			}
			newStr := reg.ReplaceAllString(record[1], "-")
			project = strings.TrimSuffix(strings.TrimPrefix(strings.ToLower(newStr), "-"), "-")
		case 1:
		default:
			csv := model.CSV{
				ID:             record[0],
				MemberName:     record[1],
				Email:          record[2],
				Role:           record[3],
				GitHubTeamRole: strings.ToLower(record[4]),
				GitHubUser:     record[5]}

			//COMMANT SKIPLINE

			if !strings.Contains(csv.ID, "#") {
				csvlist = append(csvlist, csv)
			}
		}
	}
	return nil, project, csvlist
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

func (Template) WriteCSV(teamName string, dataset []model.CSV) bool {
	file, err := os.Create(teamName + "-" + time.Now().Format("20060102150405") + ".csv")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var data = [][]string{
		{"TeamName/Project", teamName},
		{"ID (Employee)", "MemberName", "Email", "Role", "GitHub Team Role", "UserName(GitHub)"},
	}

	//prepare dataset
	for _, d := range dataset {
		data = append(data, []string{d.ID, d.MemberName, d.Email, d.Role, d.GitHubTeamRole, d.GitHubUser})
	}

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			log.Panic(err)
		}
	}

	return true
}

func (Template) ReadProjectMemberListTemplateCSV(name string) (err error, project string, projectMemberListTemplate []model.ProjectMemberListTemplate) {
	records, err := readData(name)

	if err != nil {
		log.Fatal(err)
		return err, "", nil
	}

	for i, record := range records {
		switch i {
		case 0:
			reg, err := regexp.Compile("[^A-Za-z0-9]+")
			if err != nil {
				log.Fatal(err)
			}
			project = strings.TrimSuffix(strings.TrimPrefix(strings.ToLower(reg.ReplaceAllString(record[1], "-")), "-"), "-")
		case 1:
		default:
			csv := model.ProjectMemberListTemplate{
				No:                record[0],
				Username:          record[1],
				Fullname:          record[2],
				Email:             record[3],
				Role:              record[4],
				SubscriptionOwner: record[5],
				GitHubUsername:    record[6],
				GitHubTeamRole:    strings.ToLower(record[7]),
				GitHub:            record[8],
				AzureDEV:          record[9],
				AzurePRD:          record[10],
				ELK:               record[11],
				Jumphost:          record[12],
				// Bastion:           record[13],
			}

			//COMMANT SKIPLINE

			if !strings.Contains(record[0], "#") {
				projectMemberListTemplate = append(projectMemberListTemplate, csv)
			}
		}
	}
	return nil, project, projectMemberListTemplate
}

func (Template) WriteProjectMemberListTemplateCSV(team string, header string, name string, dataset []model.ProjectMemberListTemplate) error {
	time := time.Now().Format("20060102150405")
	file, err := os.Create(name + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var data = [][]string{
		{"Project Name", team, header + time, "", "", "", "", "", "", "", "", "", "", ""},
		{"No", "Username", "Full-Name", "AIS / Postbox Email", "Role", "Subscription owner", "GitHub Username", "GitHub Role", "GitHub", "Azure DEV", "Azure PRD", "ELK", "Jumphost", "Bastion"},
	}

	//prepare dataset
	for _, d := range dataset {
		if d.Email == "" {
			d.Email = "verify email"
		}
		data = append(data, []string{d.No, d.Username, d.Fullname, d.Email, d.Role, d.SubscriptionOwner, d.GitHubUsername, d.GitHubTeamRole, d.GitHub, d.AzureDEV, d.AzurePRD, d.ELK, d.Jumphost, d.Bastion})
	}

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			log.Panic(err)
		}
	}
	return nil
}

func (Template) ReadDormantCSV(name string) (err error, dormantUsers []model.DormantUser) {
	records, err := readData(name)

	if err != nil {
		return err, nil
	}

	for i, record := range records {
		switch i {
		case 0: //header
		default:
			csv := model.DormantUser{
				CreateAt:     record[0],
				ID:           record[1],
				Login:        record[2],
				Role:         record[3],
				Suspended:    record[4],
				LastLoggedIP: record[5],
				Dormant:      record[6],
				LastActive:   record[7],
				TwoFAEnabled: record[8],
				Teams:        record[9],
				Excepted:     record[10],
			}

			//COMMANT SKIPLINE

			if !strings.Contains(record[0], "#") {
				dormantUsers = append(dormantUsers, csv)
			}
		}
	}
	return nil, dormantUsers
}

func (Template) WriteDormantCSV(name string, dataset []model.DormantUser) error {
	time := time.Now().Format("20060102150405")
	file, err := os.Create(name + "-review-" + time + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var data = [][]string{
		{"created_at", "id", "login", "role", "suspended?", "last_logged_ip", "dormant?", "last_active", "2fa_enabled?", "teams", "excepted"},
	}

	for _, d := range dataset {
		data = append(data, []string{d.CreateAt, d.ID, d.Login, d.Role, d.Suspended, d.LastLoggedIP, d.Dormant, d.LastActive, d.TwoFAEnabled, d.Teams, d.Excepted})
	}

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			log.Panic(err)
		}
	}
	return nil
}
