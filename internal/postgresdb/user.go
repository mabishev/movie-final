package postgresdb

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcokz/movie-final/internal/entity"

	"golang.org/x/crypto/bcrypt"
)

type PgxUserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(p *pgxpool.Pool) *PgxUserRepo {
	return &PgxUserRepo{pool: p}
}

func (p *PgxUserRepo) CreateUser(ctx context.Context, u entity.User) error {
	//хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = p.pool.Exec(ctx, "insert into users (email, password) values ($1, $2)", u.Email, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (p *PgxUserRepo) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var u entity.User
	err := p.pool.QueryRow(ctx, "select id, email, password from users where email = $1", email).
		Scan(&u.ID, &u.Email, &u.Password)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (p *PgxUserRepo) GetUsersByAge(ctx context.Context, minAge, maxAge int64) ([]entity.User, error) {
	rows, err := p.pool.Query(ctx, "select id, sex, dateofbirth, country, city from users where EXTRACT(YEAR FROM AGE(dateofbirth)) BETWEEN $1 AND $2",
		minAge, maxAge)
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
			return []entity.User{}, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (p *PgxUserRepo) GetUsersBySex(ctx context.Context, sex string) ([]entity.User, error) {
	rows, err := p.pool.Query(ctx, "select id, sex, dateofbirth, country, city from users where sex = $1", sex)
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
			return []entity.User{}, err // ??здесь же может быть что-то записано
		}
		users = append(users, u)
	}
	return users, nil
}

func (p *PgxUserRepo) UpdateUserInfo(ctx context.Context, u entity.User) error {
	_, err := p.pool.Exec(ctx, "update users set sex = $2, dateofbirth = $3, country = $4, city = $5 where id = $1",
		u.ID, u.Sex, u.DateOfBirth, u.Country, u.City)
	return err

}
