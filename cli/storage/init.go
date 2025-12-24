package storage

import (
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

//TODO: возможные ошибки

var MigrationsFS embed.FS

type LocalStorage struct {
	db *sql.DB
}
type product struct {
	id      int
	model   string
	company string
	price   int
}

func New() (*LocalStorage, error) {
	//TODO: перенести в конфиг для компиляции? и для тестовой БД создания
	fName := "goph_keep.sqlite"
	_, err := os.Open(fName)
	if err != nil {
		log.Printf("storage not found. err: %v", err)
		f, err := os.Create(fName)
		if err != nil {
			return nil, fmt.Errorf("could not create storage. err: %v", err)
		}
		f.Close()
		time.Sleep(1 * time.Second)
	}
	db, err := sql.Open("sqlite3", fName)
	if err != nil {
		log.Fatalf("could not open storage. err: %v", err)
	}
	// Используйте пул подключений
	db.SetMaxOpenConns(1)  // Только одно пишущее подключение
	db.SetMaxIdleConns(10) // Множество читающих

	// Включите WAL при инициализации
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return nil, fmt.Errorf("could not persist journal mode. err: %v", err)
	}
	_, err = db.Exec("PRAGMA synchronous=NORMAL;")
	if err != nil {
		return nil, fmt.Errorf("could not persist journal mode. err: %v", err)
	}
	err = migrate(fName)
	if err != nil {
		log.Fatalf("could not migrate storage. err: %v", err)
	}
	return &LocalStorage{
		db: db,
	}, nil
}

func migrate(dbURL string) error {
	//TODO: driver тоже в конфигурацию
	db, err := sql.Open("sqlite3", dbURL)
	defer func() {
		err = db.Close()
		if err != nil {
			log.Errorf("could not close db connection:%s", err)
		}
	}()
	if err != nil {
		return fmt.Errorf("could not run migration")
	}

	goose.SetBaseFS(MigrationsFS)
	if err = goose.SetDialect(string(goose.DialectSQLite3)); err != nil {
		return fmt.Errorf("could not set dialect. err:%v", err)
	}
	if err = goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("could not run migration. err: %v", err)
	}
	return nil
}
