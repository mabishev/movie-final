package postgresdb

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcokz/movie-final/internal/entity"
)

type PgxMovieRepository struct {
	pool *pgxpool.Pool
}

func NewMoviesRepository(p *pgxpool.Pool) *PgxMovieRepository {
	return &PgxMovieRepository{pool: p}
}

func (p *PgxMovieRepository) CreateMovie(ctx context.Context, m entity.Movie) error {
	_, err := p.pool.Exec(ctx, "insert into movie (name, year) values ($1, $2)", m.Name, m.Year)
	if err != nil {
		return errors.New("the movie already exists")
	}

	return nil
}

func (p *PgxMovieRepository) GetMovies(ctx context.Context) ([]entity.Movie, error) {
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

func (p *PgxMovieRepository) GetMoviesByID(ctx context.Context, id int64) (entity.Movie, error) {
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

func (p *PgxMovieRepository) UpdateMovieByID(ctx context.Context, m entity.Movie) error {
	result, err := p.pool.Exec(ctx, "update movie set name = $2, year = $3 where id = $1", m.ID, m.Name, m.Year)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return entity.ErrMovieNotFound
	}

	return nil
}

func (p *PgxMovieRepository) DeleteMovieByID(ctx context.Context, id int64) error {
	result, err := p.pool.Exec(ctx, "delete from movie where id = $1", id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return entity.ErrMovieNotFound
	}

	return nil
}
