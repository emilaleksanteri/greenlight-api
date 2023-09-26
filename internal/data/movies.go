package data

import (
	"time"

	"github.com/emilaleksanteri/greenlight-api/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateMovie(valid *validator.Validator, movie *Movie) {
	valid.Check(movie.Title != "", "title", "must be provided")
	valid.Check(len(movie.Title) <= 500, "title", "must not ne more than 500 bytes long")

	valid.Check(movie.Year != 0, "year", "must be provided")
	valid.Check(movie.Year >= 1888, "year", "must be greater than 1888")

	valid.Check(movie.Runtime != 0, "runtime", "must be provided")
	valid.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	valid.Check(movie.Genres != nil, "genres", "must be provided")
	valid.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	valid.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")

	valid.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")

}
