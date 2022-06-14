package csv

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/tarathep/ghmgr/model"
)

type Pending struct{}
type Cache struct{}

func (Template) ReadCache(name string) (err error, cache []model.Cache) {
	records, err := readData(name)

	if err != nil {
		log.Fatal(err)
		return err, nil
	}

	for i, record := range records {
		switch i {
		case 0:
		default:
			csv := model.Cache{
				No:       record[0],
				ID:       record[1],
				Username: record[2],
				Email:    record[3],
				Team:     record[4],
			}
			cache = append(cache, csv)
		}
	}
	return nil, cache
}

func (Template) WriteCache(name string, dataset []model.Cache) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var data = [][]string{
		{"No", "ID", "Username", "Email", "Team"},
	}

	//prepare dataset
	for _, d := range dataset {
		data = append(data, []string{d.No, d.ID, d.Username, d.Email, d.Team})
	}

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			log.Panic(err)
		}
	}

	return nil
}

func (Pending) WriteCache(name string, dataset []model.CachePending) error {
	os.Mkdir("cache/invited", 0755)

	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var data = [][]string{
		{"No", "ID", "Email", "Username"},
	}

	//prepare dataset
	for _, d := range dataset {
		data = append(data, []string{d.No, d.ID, d.Email, d.Username})
	}

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			log.Panic(err)
		}
	}
	return nil
}

func (Pending) ReadCache(name string) (err error, cache []model.CachePending) {
	records, err := readData(name)

	if err != nil {
		log.Fatal(err)
		return err, nil
	}

	for i, record := range records {
		switch i {
		case 0:
		default:
			csv := model.CachePending{
				No:       record[0],
				ID:       record[1],
				Email:    record[2],
				Username: record[3],
			}
			cache = append(cache, csv)
		}
	}
	return nil, cache
}
