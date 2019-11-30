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
		SELECT domain, title, ssl_grade, previous_ssl_Grade, created_at, 
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

func InsertServer(domain string, server ...Server) ([]Server, error) {
	query := `
		INSERT INTO server(address, ssl_grade, country, owner, domain)
		VALUES ($1, $2, $3, $4, $5)
	`

	for _, s := range server {
		_, err := Db.Exec(query, s.Address, s.Grade, s.Country,
			s.Owner, domain)

		if err != nil {
			return []Server{}, err
		}
	}

	return server, nil
}

func DeleteServer(ipAddress string) error {
	query := `
		DELETE FROM server WHERE ipAddress = $1
	`
	_, err := Db.Exec(query, ipAddress)
	if err != nil {
		return err
	}
	return nil
}

// Only change data that is expected to change over time on the test
func PartialUpdateSite(domain string, site Site, previousSslGrade string) error {
	var err error
	if previousSslGrade == "" {
		query := `
		UPDATE site SET (ssl_grade, previous_ssl_grade, servers_changed, updated_at)  
		= ($1, $2, $3, $4) WHERE domain = $5
	`
		_, err = Db.Exec(query, site.Grade, previousSslGrade, site.ServersChanged,
			site.UpdatedAt, domain)
	} else {
		query := `
		UPDATE site SET (ssl_grade, servers_changed, updated_at)  
		= ($1, $2, $3, $4) WHERE domain = $5
	`
		_, err = Db.Exec(query, site.Grade, site.ServersChanged,
			site.UpdatedAt, domain)
	}

	if err != nil {
		return err
	}
	return nil
}
