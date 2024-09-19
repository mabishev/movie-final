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

// func (p *PgxRatingsRepo) GetFromMovieUsersWithRating(ctx context.Context, movieid, minrating, maxrating int64) ([]entity.Rating, error) {
// 	rows, err := p.pool.Query(ctx, "select userid, rating from ratings where movieid = $1 and rating BETWEEN $2 and $3", movieid, minrating, maxrating)
// 	if err != nil {
// 		return
// 	}
// 	defer rows.Close()

// }

func (p *PgxRatingsRepo) GetMoviesWithRatingFromUser(ctx context.Context, userid, minrating, maxrating int64) ([]entity.MovieWithRating, error) {
	rows, err := p.pool.Query(ctx, "select movieid, rating from ratings where userid = $1 and rating BETWEEN $2 AND $3", userid, minrating, maxrating)
	if err != nil {
		return []entity.MovieWithRating{}, err
	}
	defer rows.Close()

	var movies []entity.MovieWithRating

	for rows.Next() {
		var m entity.MovieWithRating
		var rating int64
		err := rows.Scan(
			&rating,
			&m.Movies.ID,
			&m.Movies.Name,
			&m.Movies.Year,
		)
		if err != nil {
			return []entity.MovieWithRating{}, err
		}
		movies = append(movies, m)
	}

	if err := rows.Err(); err != nil {
		return []entity.MovieWithRating{}, err
	}

	return movies, nil
}

func (p *PgxRatingsRepo) GetUsersByRatingOfMovie(ctx context.Context, movieid, minrating, maxrating int64) ([]entity.UserWithRating, error) {
	rows, err := p.pool.Query(ctx, `
	select u.id, u.name, u.surname, u.sex, u.dateofbirth, u.country, u.city, r.rating  
	from ratings r
	JOIN users u ON r.userid = u.id
	where r.movieid = $1 AND r.rating BETWEEN $2 AND $3
	`, movieid, minrating, maxrating)
	if err != nil {
		return []entity.UserWithRating{}, err
	}
	defer rows.Close()

	var usersWithRating []entity.UserWithRating

	for rows.Next() {
		var u entity.User
		var rating int64
		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Surname,
			&u.Sex,
			&u.DateOfBirth,
			&u.Country,
			&u.City,
			&rating,
		)
		if err != nil {
			return []entity.UserWithRating{}, err
		}
		if u.ID != 0 {
			usersWithRating = append(usersWithRating, entity.UserWithRating{Users: u, Rating: rating})
		}
	}

	if err := rows.Err(); err != nil {
		return []entity.UserWithRating{}, err
	}

	return usersWithRating, err
}

func (p *PgxRatingsRepo) UpdateRating(ctx context.Context, r entity.Rating) error {
	result, err := p.pool.Exec(ctx, "insert into ratings(userid, movieid, rating) values ($1, $2, $3) ON CONFLICT (userid, movieid) DO UPDATE SET rating = EXCLUDED.rating",
		r.UserId, r.MovieID, r.Rating)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return entity.ErrMovieNotFound //??
	}

	return nil
}
