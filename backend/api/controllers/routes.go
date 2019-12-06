package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"../models"
	"../utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// For some reason the time.Since(last_update) always gives a minimun of
// four hour difference, so HoursToCheckUpdate is set to five at the
// moment to represent one hour difference.
// TODO: Fix time difference issue
const hoursToCheckUpdate = 5

type appError struct {
	Message string `json:"message"`
}

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/history", GetSiteHistory)
	router.Get("/sites/{domain}", GetSite)
	return router
}

func GetSiteHistory(w http.ResponseWriter, r *http.Request) {
	if siteMap, err := models.RetrieveLatestSites(); err != nil {
		log.Println(err)
		render.JSON(w, r, "")
	} else {
		render.JSON(w, r, siteMap)
	}
}

func GetSite(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")
	site, err := models.FetchSite(domain)
	if err != nil {
		if err == sql.ErrNoRows {
			site, err = retrieveAndSaveSite(domain)
			if err != nil {
				render.JSON(w, r, appError{Message: err.Error()})
				return
			}
			render.JSON(w, r, site)
			return
		}

		render.JSON(w, r, appError{Message: err.Error()})
		return
	}

	timeDiff := time.Since(site.UpdatedAt)
	if timeDiff.Hours() > hoursToCheckUpdate {
		site, err = checkForUpdate(site)
		if err != nil {
			render.JSON(w, r, appError{Message: err.Error()})
			return
		}
	}

	site.Servers, _ = models.FetchServers(domain)
	render.JSON(w, r, site)
}

func performServerUpdate(site models.Site) (models.Site, error) {
	updatedSite, status := utils.GetWebsiteData(site.Domain)

	if status != utils.InProgressMsg {
		updatedSite.ServersChanged = true
		updatedSite.PreviousGrade = site.Grade

		// Delete and add new the servers
		for _, s := range site.Servers {
			models.DeleteServer(s.Address)
		}
		models.InsertServer(site.Domain, updatedSite.Servers...)

		err := models.PartialUpdateSite(updatedSite, site.Domain)
		if err != nil {
			return models.Site{}, err
		}
		return updatedSite, nil
	}
	return models.Site{}, errors.New(utils.InProgressMsg)
}

func checkForUpdate(site models.Site) (models.Site, error) {
	updatedSite := site

	updated := utils.HasServersUpdated(site.Domain, site.Servers)

	if updated {
		var err error
		updatedSite, err = performServerUpdate(site)
		if err != nil {
			return models.Site{}, err
		}
	} else if site.ServersChanged == true {
		site.ServersChanged = false
		models.PartialUpdateSite(site, "")
		updatedSite = site
	}
	return updatedSite, nil
}

func retrieveAndSaveSite(domain string) (models.Site, error) {
	site, status := utils.GetWebsiteData(domain)
	if status == utils.ReadyMsg {
		if _, err := models.InsertSite(site); err != nil {
			return models.Site{}, err
		}
		if _, err := models.InsertServer(site.Domain, site.Servers...); err != nil {
			return models.Site{}, err
		}
		return site, nil
	} else if site.IsDown {
		return site, nil
	}
	return site, errors.New(utils.InProgressMsg)
}
