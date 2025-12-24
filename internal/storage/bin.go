package storage

import (
	"context"
	"github.com/randallmlough/pgxscan"
)

type Bin struct {
	Bin  string `json:"bin"`
	Meta string `json:"meta"`
}

func (p *PGStorage) StoreBin(ctx context.Context, userID int, dataName, bin, meta string) error {
	sqlQuery := `INSERT INTO bin (user_id, name, bin, meta) VALUES($1,$2,$3,$4) 
				 ON CONFLICT(name) DO UPDATE SET bin=excluded.bin, meta=excluded.meta`
	_, err := p.Connection.Exec(ctx, sqlQuery, userID, dataName, bin, meta)
	if err != nil {
		return err
	}
	return nil
}

func (p *PGStorage) GetBin(ctx context.Context, userID int, dataName string) (*Bin, error) {
	var resp Bin
	sqlQuery := `SELECT bin, meta FROM bin WHERE user_id=$1 AND name=$2`
	row := p.Connection.QueryRow(ctx, sqlQuery, userID, dataName)
	err := pgxscan.NewScanner(row).Scan(&resp.Bin, &resp.Meta)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
