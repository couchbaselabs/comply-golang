package main

import (
	"time"

	"github.com/couchbase/gocb"
)

type Company struct {
	Type      string `json:"_type"`
	ID        string `json:"_id"`
	CreatedOn string `json:"createdON"`
	Name      string `json:"name"`
	Address   struct {
		Street  string `json:"street"`
		City    string `json:"city"`
		State   string `json:"state"`
		Zip     int    `json:"zip"`
		Country string `json:"country"`
	} `json:"address"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
	Active  bool   `json:"active"`
}

type SessionCompany struct {
	Company Company
}

func (c *SessionCompany) Create() (*Company, error) {
	c.Company.Type = "Company"
	c.Company.ID = c.Company.Website
	c.Company.CreatedOn = time.Now().Format(time.RFC3339)
	c.Company.Active = true

	_, err := bucket.Upsert(c.Company.Website, c.Company, 0)
	if err != nil {
		return nil, err
	}
	return &c.Company, nil
}

func (c *SessionCompany) Retrieve(id string) (*Company, error) {
	_, err := bucket.Get(id, &c.Company)
	if err != nil {
		return nil, err
	}
	return &c.Company, nil
}

func (c *SessionCompany) RetrieveAll() ([]Company, error) {
	myQuery := gocb.NewN1qlQuery("SELECT * FROM `comply` " +
		"WHERE _type = 'Company'")
	rows, err := bucket.ExecuteN1qlQuery(myQuery, nil)
	if err != nil {
		return nil, err
	}

	type wrapCompany struct {
		Company Company `json:"comply"`
	}
	var row wrapCompany
	var curCompanies []Company

	for rows.Next(&row) {
		curCompanies = append(curCompanies, row.Company)
	}
	return curCompanies, nil
}
