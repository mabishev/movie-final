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

// func (p *PgxRatingsRepo) GetMovieByRating(ctx context.Context, minrating, maxrating int64) ([]entity.Movie, error) {
// 	p.pool.Query(ctx, "select movieid")
// }

type UserWithRating struct {
	User   entity.User
	Rating int64
}

func (p *PgxRatingsRepo) GetUsersByRatingOfMovie(ctx context.Context, movieid, minrating, maxrating int64) ([]UserWithRating, error) {
	rows, err := p.pool.Query(ctx, `
	select r.rating, u.id, u.name, u.surname, u.sex, u.dateofbirth, u.country, u.city 
	from ratings r
	JOIN users u ON r.userid = u.id
	where r.movieid = $1 AND r.rating BETWEEN $2 AND $3
	`, movieid, minrating, maxrating)
	if err != nil {
		return []UserWithRating{}, err
	}
	defer rows.Close()

	var usersWithRating []UserWithRating

	for rows.Next() {
		var u entity.User
		var rating int64
		err := rows.Scan(
			&rating,
			&u.ID,
			&u.Name,
			&u.Surname,
			&u.Sex,
			&u.DateOfBirth,
			&u.Country,
			&u.City,
		)
		if err != nil {
			return nil, err
		}
		usersWithRating = append(usersWithRating, UserWithRating{User: u, Rating: rating})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usersWithRating, err
}

func (p *PgxRatingsRepo) UpdateRating(ctx context.Context, r entity.Rating) error {
	_, err := p.pool.Exec(ctx, "insert into ratings(userid, movieid, rating) values ($1, $2, $3) ON CONFLICT (userid, movieid) DO UPDATE SET rating = EXCLUDED.rating",
		r.UserId, r.MovieID, r.Rating)
	if err != nil {
		return err
	}
	return nil
}
