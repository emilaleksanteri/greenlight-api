package data

import (
	"database/sql"
	"fmt"
	"time"

	"errors"

	"context"

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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
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
	query := `
		update movies
		set title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
		where id = $5 and version = $6
		returning version
	`
	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		delete from movies
		where id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {
	query := fmt.Sprintf(`
		select count(*) over(), id, created_at, title, year, runtime, genres, version
		from movies
		where (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) or $1 = '')
		and (genres @> $2 or $2 = '{}')
		order by %s %s, id asc
		limit $3 offset $4`,
		filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{title, pq.Array(genres), filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()
	movies := []*Movie{}
	totalRecords := 0

	for rows.Next() {
		var movie Movie

		err := rows.Scan(
			&totalRecords,
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		movies = append(movies, &movie)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return movies, metadata, nil
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
