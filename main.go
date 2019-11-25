package main

import (
	"database/sql"
	"log"
	"net/http"

	"./api"
	"./api/controllers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Route("/", func(r chi.Router) {
		r.Mount("/servers", controllers.Routes())
	})

	return router
}

var db *sql.DB

func GetDb() *sql.DB {
	return db
}

func main() {
	router := Routes()

	walkFunc := func(method string, route string, handler http.Handler, middleware ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}

	db = api.InitDb()

	log.Fatal(http.ListenAndServe(":8000", router))
}
