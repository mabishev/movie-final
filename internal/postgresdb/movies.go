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

func (p *PgxMovieRepository) CreateMovie(ctx context.Context, e entity.Movie) error {
	_, err := p.pool.Exec(ctx, "insert into movie (name, year) values ($1, $2)", e.Name, e.Year)
	if err != nil {
		return errors.New("the movie already exists")
	}
	return err
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
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

func (p *PgxMovieRepository) GetMoviesByID(ctx context.Context, id int64) (entity.Movie, error) {
	row := p.pool.QueryRow(ctx, "select id, name, year from movie where id = $1", id)

	var e entity.Movie
	err := row.Scan(&e.ID, &e.Name, &e.Year)
	if errors.Is(err, pgx.ErrNoRows) { //???
		return entity.Movie{}, entity.ErrNotFound
	}
	if err != nil {
		return entity.Movie{}, err
	}

	return e, nil
}

func (p *PgxMovieRepository) UpdateMovie(ctx context.Context, e entity.Movie) error {
	result, err := p.pool.Exec(ctx, "update movie set name = $2, year = $3 where id = $1", e.ID, e.Name, e.Year)
	if err != nil {
		return err
	}
	// Возвращаем ErrNotFound, если запись не была найдена и обновлена
	if result.RowsAffected() == 0 {
		return entity.ErrNotFound
	}
	return nil
}

func (p *PgxMovieRepository) DeleteMovieByID(ctx context.Context, id int64) error {
	_, err := p.pool.Exec(ctx, "delete from movie where id = $1", id)
	return err
}
