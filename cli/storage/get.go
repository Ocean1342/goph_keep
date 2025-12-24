package storage

import (
	"context"
	"fmt"
	"goph_keeper/cli/entity"
	"goph_keeper/pkg/auth"
	"time"
)

func (ls LocalStorage) GetLocalUserByLogin(ctx context.Context, login string) (*entity.User, error) {
	var user entity.User
	sqlQuery := `SELECT * FROM users WHERE login=$1`
	row := ls.db.QueryRowContext(ctx, sqlQuery, login)
	var token *auth.Token
	var authDT string
	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Phrase, &token, &user.IsAuthenticated, &authDT)
	if token != nil {
		user.Token = string(*token)
	}
	if authDT != "" {
		parse, err := time.Parse(time.DateTime, authDT)
		if err != nil {
			return nil, err
		}
		user.AuthDt = parse
	}
	if err != nil {
		return nil, fmt.Errorf("get user by login failed: %v", err)
	}
	return &user, nil
}

func (ls LocalStorage) GetLogoPass(ctx context.Context, userID int, name string) (logo, pass, meta string, err error) {
	sqlQuery := `select creeds.password, creeds.login, creeds.meta
	from catalog
	join creeds on catalog.data_id = creeds.id
	where catalog.user_id = $1
	and catalog.data_name = $2`
	row := ls.db.QueryRowContext(ctx, sqlQuery, userID, name)
	err = row.Scan(&logo, &pass, &meta)
	if err != nil {
		return "", "", "", fmt.Errorf("get user data failed: %v", err)
	}
	return logo, pass, meta, nil
}

func (ls LocalStorage) GetBin(ctx context.Context, userID int, name string) (bin, meta string, err error) {
	sqlQuery := `select bin.bin, bin.meta
	from catalog
	join bin on catalog.data_id = bin.id
	where catalog.user_id = $1
	and catalog.data_name = $2`
	row := ls.db.QueryRowContext(ctx, sqlQuery, userID, name)
	err = row.Scan(&bin, &meta)
	if err != nil {
		return "", "", fmt.Errorf("get user data failed: %v", err)
	}
	return bin, meta, nil
}
