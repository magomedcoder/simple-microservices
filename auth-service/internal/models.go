package internal

import (
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{
		User: User{},
	}
}

type Models struct {
	User User
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Name      string    `json:"name,omitempty"`	
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var user User
	row := db.QueryRowContext(ctx, `SELECT id, email, password, name, created_at FROM users WHERE email = $1`, email)
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
