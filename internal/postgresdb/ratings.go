package postgresdb

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxRatingsRepo struct {
	pool *pgxpool.Pool
}

func NewRatingsRepo(p *pgxpool.Pool) *PgxRatingsRepo {
	return &PgxRatingsRepo{pool: p}
}

func (p *PgxRatingsRepo) AddRating(ctx context.Context, userid, movieid, rating string) error {
	_, err := p.pool.Exec(ctx, "insert into ratings (rating) values = $3 where userid = $1, movieid = $2",
		userid, movieid, rating)
	if err != nil {
		return err
	}
	return nil
}
