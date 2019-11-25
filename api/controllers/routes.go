package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"../models"
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

	if err != nil && err == sql.ErrNoRows {
		log.Fatal(err)
	} else {
		log.Fatal(err)
	}
	render.JSON(w, r, site)
}
