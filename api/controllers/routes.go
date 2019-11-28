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

	} else if err != nil {
		log.Println(err)
	}
	render.JSON(w, r, site)
}
