package models

import (
	"database/sql"
	"time"
)

var Db *sql.DB

type Site struct {
	Domain         string    `json:"domain"`
	Title          string    `json:"title"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updateddAt"`
	Grade          string    `json:"grade"`
	PreviousGrade  string    `json:"previousGrade"`
	Logo           string    `json:"logo"`
	IsDown         bool      `json:"isDown"`
	ServersChanged bool      `json:"serversChanged"`

	Servers []Server `json:"servers"`
}

type Server struct {
	Address string `json:"address"`
	Grade   string `json:"sslGrade"`
	Country string `json:"country"`
	Owner   string `json:"owner"`
}

func FetchSite(domain string) (Site, error) {
	query := `
		SELECT domain, title, ssl_grade, previous_ssl_grade, created_at, 
		       updated_at, logo, is_down
		FROM site 
		WHERE domain = $1
	`
	site := Site{}
	if err := Db.QueryRow(query, domain).Scan(&site.Domain, &site.Title,
		&site.Grade, &site.PreviousGrade, &site.CreatedAt, &site.UpdatedAt,
		&site.Logo, &site.IsDown); err != nil {
		return Site{}, err
	}

	return site, nil
}

func InsertSite(site Site) (Site, error) {
	query := `
		INSERT INTO site(domain, title, ssl_grade, previous_ssl_grade, 
		                 created_at, updated_at, logo, is_down)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := Db.Exec(query, site.Domain, site.Title, site.Grade,
		site.PreviousGrade, site.CreatedAt, site.UpdatedAt,
		site.Logo, site.IsDown)

	if err != nil {
		return Site{}, err
	}

	return site, nil
}

func FetchServers(domain string) ([]Server, error) {
	servers := []Server{}
	query := `
		SELECT address, ssl_grade, country, owner
		FROM server
		WHERE domain = $1
	`
	rows, err := Db.Query(query, domain)
	if err != nil {
		return []Server{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var s Server
		err := rows.Scan(&s.Address, &s.Grade, &s.Country, &s.Owner)
		if err != nil {
			return []Server{}, err
		}
		servers = append(servers, s)
	}

	return servers, nil
}

func InsertServer(server Server, domain string) (Server, error) {
	query := `
		INSERT INTO server(address, ssl_grade, country, owner, domain)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := Db.Exec(query, server.Address, server.Grade, server.Country,
		server.Owner, domain)

	if err != nil {
		return Server{}, err
	}

	return server, nil
}
