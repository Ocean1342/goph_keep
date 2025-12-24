package storage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"goph_keeper/pkg/auth"
	"time"
)

type DataType string

const (
	LogoPass        DataType = "logo_pass"
	ArbitraryText   DataType = "arbitrary_text"
	ArbitraryBinary DataType = "arbitrary_binary"
	CreditCard      DataType = "credit_card"
)

func (ls LocalStorage) StoreToken(ctx context.Context, login string, token auth.Token) error {
	tx, err := ls.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx err: %w", err)
	}
	sqlQuery := `UPDATE users
				 SET token=$1, is_authenticated=$2, auth_dt=$3
				 WHERE login = $4`
	_, err = tx.Exec(sqlQuery, string(token), 1, time.Now().UTC().Format("2006-01-02 15:04:05"), login)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return fmt.Errorf("rollback tx err: %w", err)
		}
		return fmt.Errorf("could not insert user.err: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit tx err: %w", err)
	}
	return nil
}

// StoreInitialUserData - сохранение данных пользователя в локальном хранилище.
// может использоваться когда недоступен сервер
func (ls LocalStorage) StoreInitialUserData(ctx context.Context, login, pass, phrase string) error {
	tx, err := ls.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx err: %w", err)
	}
	sqlQuery := `INSERT INTO users (login,password, phrase) VALUES($1,$2,$3) 
				 ON CONFLICT DO NOTHING 
                 RETURNING id, login, password, phrase`
	_, err = tx.Exec(sqlQuery, login, pass, phrase)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return fmt.Errorf("rollback tx err: %w", err)
		}
		return fmt.Errorf("could not insert user.err: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit tx err: %w", err)
	}
	return nil
}

func (ls LocalStorage) StoreLogoPass(ctx context.Context, userID int, name, storeLogo, storePass, meta string) error {
	tx, err := ls.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx err: %w", err)
	}
	dataID := uuid.New()
	sqlQuery := `INSERT INTO catalog (user_id, data_name, data_id, data_type, is_synced) VALUES($1,$2,$3,$4,$5) 
				 ON CONFLICT DO NOTHING 
                 `
	_, err = tx.Exec(sqlQuery, userID, name, dataID.String(), string(LogoPass), 0)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return fmt.Errorf("rollback tx err: %w", err)
		}
		return fmt.Errorf("could not insert into catalog. err: %w", err)
	}

	sqlQuery2 := `INSERT INTO creeds (id, login, password, meta) VALUES($1,$2,$3,$4) 
                  ON CONFLICT(id) DO UPDATE
                  SET login=excluded.login,
                      password=excluded.password,
                      meta=excluded.meta`
	_, err = tx.Exec(sqlQuery2, dataID, storeLogo, storePass, meta)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return fmt.Errorf("rollback tx err: %w", err)
		}
		return fmt.Errorf("could not insert into creds. err: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit tx err: %w", err)
	}
	return nil
}

func (ls LocalStorage) StoreBin(ctx context.Context, userID int, name, storeBin, meta string) error {
	tx, err := ls.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx err: %w", err)
	}
	dataID := uuid.New()
	sqlQuery := `INSERT INTO catalog (user_id, data_name, data_id, data_type, is_synced) VALUES($1,$2,$3,$4,$5) 
				 ON CONFLICT DO NOTHING 
                 `
	_, err = tx.Exec(sqlQuery, userID, name, dataID.String(), string(ArbitraryBinary), 0)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return fmt.Errorf("rollback tx err: %w", err)
		}
		return fmt.Errorf("could not insert into catalog. err: %w", err)
	}

	sqlQuery2 := `INSERT INTO bin (id, bin, meta) VALUES($1,$2,$3) 
                  ON CONFLICT(id) DO UPDATE
                  SET bin = excluded.bin,
                      meta=excluded.meta`
	_, err = tx.Exec(sqlQuery2, dataID, storeBin, meta)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return fmt.Errorf("rollback tx err: %w", err)
		}
		return fmt.Errorf("could not insert into creds. err: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit tx err: %w", err)
	}
	return nil
}
