package storage

import (
	"context"
	"fmt"
)

type CatalogItem struct {
	UserID    int    `db:"user_id"`
	UserLogin string `db:"user_name"`
	UserToken string `db:"token"`
	DataName  string `db:"data_name"`
	DataID    string `db:"data_id"`
	DataType  string `db:"data_type"`
	IsSynced  bool   `db:"is_synced"`
}

func (ls LocalStorage) GetNotSyncedItems(ctx context.Context) ([]CatalogItem, error) {
	sqlQuery := `select users.login, catalog.user_id, users.token, catalog.data_name, catalog.data_id, catalog.data_type, catalog.is_synced
				from catalog
				join users on catalog.user_id = users.id
				where catalog.is_synced = 0`
	res, err := ls.db.QueryContext(ctx, sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("catalog item query failed: %w", err)
	}
	var items []CatalogItem
	for res.Next() {
		var item CatalogItem
		err = res.Scan(&item.UserLogin, &item.UserID, &item.UserToken, &item.DataName, &item.DataID, &item.DataType, &item.IsSynced)
		if err != nil {
			return nil, fmt.Errorf("catalog item scan failed: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (ls LocalStorage) UpdateCatalogItem(ctx context.Context, userID int, dataID string) error {
	sqlQuery := `update catalog set is_synced = 1 where user_id = $1 and data_id = $2`
	_, err := ls.db.ExecContext(ctx, sqlQuery, userID, dataID)
	if err != nil {
		return fmt.Errorf("catalog item update failed: %w", err)
	}
	return nil
}
