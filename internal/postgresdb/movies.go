package postgresdb

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcokz/movie-final/internal/movie"
)

type PgxMovieRepository struct {
	pool *pgxpool.Pool
}

func NewMoviesRepository(p *pgxpool.Pool) *PgxMovieRepository {
	return &PgxMovieRepository{pool: p}
}

func (p *PgxMovieRepository) CreateMovie(ctx context.Context, m movie.Movie) error {
	_, err := p.pool.Exec(ctx, "insert into movie (name, year) values ($1, $2)", m.Name, m.Year)
	if err != nil {
		return errors.New("the movie already exists")
	}
	return err
}

func (p *PgxMovieRepository) GetMovies(ctx context.Context) []movie.Movie {
	rows, err := p.pool.Query(ctx, "select id, name, year from movie")
	if err != nil {
		return nil
	}
	var movies []movie.Movie

	for rows.Next() {
		var m movie.Movie
		err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Year,
		)
		if err != nil {
			return nil
		}
		movies = append(movies, m)
	}
	return movies
}

func (p *PgxMovieRepository) GetMoviesByID(ctx context.Context, id int64) (movie.Movie, error) {
	row := p.pool.QueryRow(ctx, "select id, name, year from movie where id = $1", id)

	var m movie.Movie
	err := row.Scan(&m.ID, &m.Name, &m.Year)
	if errors.Is(err, pgx.ErrNoRows) { //???
		return movie.Movie{}, movie.ErrNotFound
	}
	if err != nil {
		return movie.Movie{}, err
	}

	return m, nil
}

func (p *PgxMovieRepository) UpdateMovie(ctx context.Context, m movie.Movie) error {
	result, err := p.pool.Exec(ctx, "update movie set name = $2, year = $3 where id = $1", m.ID, m.Name, m.Year)
	if err != nil {
		return err
	}
	// если запрос не затронет ни одной строки, то вернет ошибку
	if result.RowsAffected() == 0 {
		return movie.ErrNotFound
	}
	return nil
}

func (p *PgxMovieRepository) DeleteMovieByID(ctx context.Context, id int64) error {
	_, err := p.pool.Exec(ctx, "delete from movie where id = $1", id)
	return err
}
