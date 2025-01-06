package postgresdb

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcokz/movie-final/internal/entity"
)

type PgxMoviesRepo struct {
	pool *pgxpool.Pool
}

func NewMoviesRepo(p *pgxpool.Pool) *PgxMoviesRepo {
	return &PgxMoviesRepo{pool: p}
}

func (p *PgxMoviesRepo) CreateMovie(ctx context.Context, m entity.Movie) error {
	_, err := p.pool.Exec(ctx, "insert into movie (name, year, description) values ($1, $2, $3)", m.Name, m.Year, m.Description)
	if err != nil {
		return errors.New("the movie already exists")
	}

	return nil
}

func (p *PgxMoviesRepo) GetMovies(ctx context.Context) ([]entity.Movie, error) {
	rows, err := p.pool.Query(ctx, "select id, name, year from movie")
	if err != nil {
		return []entity.Movie{}, err
	}
	defer rows.Close()

	var movies []entity.Movie

	for rows.Next() {
		var m entity.Movie
		err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Year,
		)
		if err != nil {
			return []entity.Movie{}, err
		}
		movies = append(movies, m)
	}

	if err := rows.Err(); err != nil {
		return []entity.Movie{}, err
	}

	return movies, nil
}

func (p *PgxMoviesRepo) GetMoviesByID(ctx context.Context, id int64) (entity.Movie, error) {
	var e entity.Movie

	err := p.pool.QueryRow(ctx, "select id, name, year from movie where id = $1", id).
		Scan(&e.ID, &e.Name, &e.Year)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Movie{}, entity.ErrMovieNotFound
		}
		return entity.Movie{}, err
	}

	return e, nil
}

func (p *PgxMoviesRepo) UpdateMovieByID(ctx context.Context, m entity.Movie) error {
	result, err := p.pool.Exec(ctx, "update movie set name = $2, year = $3 where id = $1", m.ID, m.Name, m.Year)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return entity.ErrMovieNotFound
	}

	return nil
}

func (p *PgxMoviesRepo) DeleteMovieByID(ctx context.Context, id int64) error {
	result, err := p.pool.Exec(ctx, "delete from movie where id = $1", id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return entity.ErrMovieNotFound
	}

	return nil
}
