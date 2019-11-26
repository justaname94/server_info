package models

import (
	"database/sql"
	"log"
	"time"
)

var Db *sql.DB

type Site struct {
	Domain           string    `json:"domain"`
	Title            string    `json:"title"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updateddAt"`
	SSLGrade         string    `json:"sslGrade"`
	PreviousSSLGrade string    `json:"previousSSLGrade"`
	Logo             string    `json:"logo"`
	IsDown           bool      `json:"isDown"`
	ServersChanged   bool      `json:"serversChanged"`
}

type Server struct {
	Address      string    `json:"address"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updateddAt"`
	SSLGrade     string    `json:"sslGrade"`
	Country      string    `json:"country"`
	Owner        string    `json:"owner"`
	Organization string    `json:"organization"`
}

func FetchSite(domain string) (Site, error) {
	query := `
		SELECT domain, title, ssl_grade, previous_ssl_grade, logo, is_down
		FROM site 
		WHERE domain = $1
	`
	site := Site{}
	if err := Db.QueryRow(query, domain).Scan(&site.Domain, &site.Title, &site.SSLGrade, &site.PreviousSSLGrade, &site.Logo, &site.IsDown); err != nil {
		return Site{}, err
	}

	return site, nil
}

func InsertSite(domain string) Site {
	query := `
		INSERT INTO site(domain, title, ssl_grade, previous_ssl_grade, logo, is_down)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	site := Site{}
	_, err := Db.Exec(query, site.Domain, site.Title, site.SSLGrade, site.PreviousSSLGrade, site.Logo, site.IsDown)

	if err != nil {
		log.Fatal(err)
	}

	return site
}
