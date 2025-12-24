package storage

import (
	"context"
	"fmt"
	"github.com/randallmlough/pgxscan"
)

type User struct {
	ID       int    `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
	Phrase   string `db:"phrase"`
}

func (p *PGStorage) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	var user User
	sqlQuery := "SELECT id, login, password FROM users WHERE login = $1"
	row := p.Connection.QueryRow(ctx, sqlQuery, login)
	err := pgxscan.NewScanner(row).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("scan rows error. err: %v", err)
	}
	return &user, nil
}

func (p *PGStorage) CreateUser(ctx context.Context, login, password, phrase string) (*User, error) {
	var user User
	sqlQuery := `INSERT INTO users (login,password, phrase) VALUES($1,$2,$3) RETURNING id, login, password, phrase`
	row := p.Connection.QueryRow(ctx, sqlQuery, login, password, phrase)
	err := pgxscan.NewScanner(row).Scan(&user.ID, &user.Login, &user.Password, &user.Phrase)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
