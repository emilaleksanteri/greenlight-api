package main

import (
	"fmt"
	"github.com/emilaleksanteri/greenlight-api/internal/data"
	"github.com/emilaleksanteri/greenlight-api/internal/validator"
	"net/http"
	"time"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input) // mf wants a non nill pointer so init input 1st
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	valid := validator.New()
	
	data.ValidateMovie(valid, &data.Movie{
		Title: input.Title,
		Year: input.Year,
		Runtime: input.Runtime,
		Genres: input.Genres,
	})
	
	if !valid.Valid() {
		app.failedValidationResponse(w, r, valid.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "The Grand Hotel Budapest",
		Runtime:   100,
		Genres:    []string{"comedy", "drama"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
