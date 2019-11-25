package controllers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Sample struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{id}", SampleRoute)
	return router
}

func SampleRoute(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	samples := Sample{
		Slug:  id,
		Title: "Hello World",
		Body:  "A simple hello world",
	}
	render.JSON(w, r, samples)
}
