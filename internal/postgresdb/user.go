package postgresdb

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = p.pool.Exec(ctx, "insert into users (email, password) values ($1, $2)", u.Email, hashedPassword)
	if err != nil {
		return errors.New("the user already exists")
	}

	return nil
}

func (p *PgxUserRepo) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var u entity.User
	err := p.pool.QueryRow(ctx, "select id, email, password from users where email = $1", email).
		Scan(&u.ID, &u.Email, &u.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, entity.ErrUserNotFound
		}
		return entity.User{}, err
	}

	return u, nil
}

func (p *PgxUserRepo) GetUserByAge(ctx context.Context, minAge, maxAge int64) ([]entity.User, error) {
	rows, err := p.pool.Query(ctx, "select id, name, surname, sex, dateofbirth, country, city from users where EXTRACT(YEAR FROM AGE(dateofbirth)) BETWEEN $1 AND $2",
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
			&u.Name,
			&u.Surname,
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

	if err := rows.Err(); err != nil {
		return []entity.User{}, err
	}
	return users, nil
}

func (p *PgxUserRepo) GetUserByCountry(ctx context.Context, country string) ([]entity.User, error) {
	rows, err := p.pool.Query(ctx, "select id, name, surname, sex, dateofbirth, country, city from users where country = $1", country)
	if err != nil {
		return []entity.User{}, err
	}
	defer rows.Close()

	var users []entity.User

	for rows.Next() {
		var u entity.User
		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Surname,
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

	if err := rows.Err(); err != nil {
		return []entity.User{}, err
	}

	return users, nil
}

func (p *PgxUserRepo) GetUserByCity(ctx context.Context, city string) ([]entity.User, error) {
	rows, err := p.pool.Query(ctx, "select id, name, surname, sex, dateofbirth, country, city from users where city = $1", city)
	if err != nil {
		return []entity.User{}, err
	}
	defer rows.Close()

	var users []entity.User

	for rows.Next() {
		var u entity.User
		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Surname,
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

	if err := rows.Err(); err != nil {
		return []entity.User{}, err
	}

	return users, nil
}

func (p *PgxUserRepo) GetUserBySex(ctx context.Context, sex string) ([]entity.User, error) {
	rows, err := p.pool.Query(ctx, "select id, name, surname, sex, dateofbirth, country, city from users where sex = $1", sex)
	if err != nil {
		return []entity.User{}, err
	}
	defer rows.Close()

	var users []entity.User

	for rows.Next() {
		var u entity.User
		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Surname,
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

	if err := rows.Err(); err != nil {
		return []entity.User{}, err
	}

	return users, nil
}

func (p *PgxUserRepo) UpdateUserInfo(ctx context.Context, u entity.User) error {
	result, err := p.pool.Exec(ctx, "update users set name = $2, surname = $3, sex = $4, dateofbirth = $5, country = $6, city = $7 where id = $1",
		u.ID, u.Name, u.Surname, u.Sex, u.DateOfBirth, u.Country, u.City)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return entity.ErrUserNotFound
	}

	return nil
}
