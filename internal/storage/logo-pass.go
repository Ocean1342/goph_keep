package storage

import (
	"context"
	"github.com/randallmlough/pgxscan"
)

type LogoPass struct {
	Login string `json:"user_login"`
	Pass  string `json:"user_pass"`
	Meta  string `json:"meta"`
}

func (p *PGStorage) StoreLogoPass(ctx context.Context, userID int, dataName, sLogo, sPass, meta string) error {
	sqlQuery := `INSERT INTO creeds (user_id, name, login, password, meta) VALUES($1,$2,$3,$4,$5) 
				 ON CONFLICT(name) DO UPDATE SET login=excluded.login, password=excluded.password, meta=excluded.meta`
	_, err := p.Connection.Exec(ctx, sqlQuery, userID, dataName, sLogo, sPass, meta)
	if err != nil {
		return err
	}
	return nil
}

func (p *PGStorage) GetLogoPass(ctx context.Context, userID int, dataName string) (*LogoPass, error) {
	var resp LogoPass
	sqlQuery := `SELECT login, password, meta FROM creeds WHERE user_id=$1 AND name=$2`
	row := p.Connection.QueryRow(ctx, sqlQuery, userID, dataName)
	err := pgxscan.NewScanner(row).Scan(&resp.Login, &resp.Pass, &resp.Meta)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
