package data

import (
	"database/sql"
	"time"

	"errors"
	"github.com/emilaleksanteri/greenlight-api/internal/validator"
	"github.com/lib/pq"
)

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(movie *Movie) error {
	query := `
		insert into movies (title, year, runtime, genres)
		values ($1, $2, $3, $4)
		returning id, created_at, version
	`
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		select id, created_at, title, year, runtime, genres, version
		from movies
		where id = $1
	`
	var movie Movie

	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {
	return nil
}

func (m MovieModel) Delete(id int64) error {
	return nil
}

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
