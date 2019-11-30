package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"../models"
	"../utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

const HoursToCheckUpdate = 5

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{domain}", GetSite)
	return router
}

func GetSite(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")
	site, err := models.FetchSite(domain)
	if err == sql.ErrNoRows {
		site = utils.GetWebsiteData(domain)
		site.CreatedAt = time.Now()
		site.UpdatedAt = time.Now()
		_, err := models.InsertSite(site)

		if err != nil {
			log.Println(err)
		}

		_, err = models.InsertServer(site.Domain, site.Servers...)
		if err != nil {
			log.Println(err)
		}

	} else if err != nil {
		log.Println(err)
	} else {
		timeDiff := time.Since(site.UpdatedAt)
		if err != nil {
			log.Println(err)
		}
		servers, _ := models.FetchServers(domain)
		var updated bool
		// For some reason the time.Since(last_update) always gives a minimun of
		// four hour difference, so HoursToCheckUpdate is set to five at the
		// moment to represent one hour difference.
		// TODO: Fix time difference issue
		if timeDiff.Hours() > HoursToCheckUpdate {
			updated = utils.HasServersUpdated(domain, servers)
		}

		if updated {
			updatedSite := utils.GetWebsiteData(domain)
			updatedSite.ServersChanged = true
			updatedSite.PreviousGrade = site.Grade
			updatedSite.CreatedAt = site.CreatedAt
			updatedSite.UpdatedAt = time.Now()

			// Delete and add new the servers
			for _, s := range site.Servers {
				models.DeleteServer(s.Address)
			}
			models.InsertServer(domain, updatedSite.Servers...)

			err = models.PartialUpdateSite(domain, updatedSite, site.Grade)
			if err != nil {
				log.Println(err)
			}
			site = updatedSite
		} else {
			site.Servers = servers
			site.ServersChanged = false
			models.PartialUpdateSite(domain, site, "")
			if err != nil {
				log.Println(err)
			}
		}
	}
	if site.Logo != "" {
		// TODO: Automate host URI to static content
		site.Logo = r.Host + site.Logo
	}
	render.JSON(w, r, site)
}
