package postgresdb

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcokz/movie-final/internal/entity"
)

type PgxRatingsRepo struct {
	pool *pgxpool.Pool
}

func NewRatingsRepo(p *pgxpool.Pool) *PgxRatingsRepo {
	return &PgxRatingsRepo{pool: p}
}

func (p *PgxRatingsRepo) GetUsersByRating(ctx context.Context, movieid, minrating, maxrating int64) ([]entity.User, error) {
	rows, err := p.pool.Query(ctx, "select userid from ratings where movieid = $1 AND rating BETWEEN $2 AND $3", movieid, minrating, maxrating)
	if err != nil {
		return []entity.User{}, err
	}
	defer rows.Close()

	var users []entity.User

	for rows.Next() {
		var u entity.User
		err := rows.Scan(
			&u.ID,
			&u.Sex,
			&u.DateOfBirth,
			&u.Country,
			&u.City,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, err
}

func (p *PgxRatingsRepo) UpdateRating(ctx context.Context, r entity.Rating) error {
	_, err := p.pool.Exec(ctx, "insert into ratings (rating) values = $3 where userid = $1, movieid = $2",
		r.UserId, r.MovieID, r.Rating)
	if err != nil {
		return err
	}
	return nil
}
