package postgresdb

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcokz/movie-final/internal/users"
	"golang.org/x/crypto/bcrypt"
)

type PgxUserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(p *pgxpool.Pool) *PgxUserRepo {
	return &PgxUserRepo{pool: p}
}

func (p *PgxUserRepo) CreateUser(ctx context.Context, u users.User) error {
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

func (p *PgxUserRepo) GetUserByEmail(ctx context.Context, email string) (users.User, error) {
	// Получение юзера из базы данных
	var u users.User
	err := p.pool.QueryRow(ctx, "select id, email, password from users where email = $1", email).
		Scan(&u.ID, &u.Email, &u.Password) //query for login or email??
	if err != nil {
		return u, err
	}

	return u, nil
}

func (p *PgxUserRepo) UpdateUserInfo(ctx context.Context, u users.User) error {
	_, err := p.pool.Exec(ctx, "update users set sex = $2, dateofbirth = $3, country = $4, city = $5 where id = $1",
		u.ID, u.Sex, u.DateOfBirth, u.Country, u.City)
	return err

}
