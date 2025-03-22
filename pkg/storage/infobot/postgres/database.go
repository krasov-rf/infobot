package infobotdb_pg

import (
	"context"
	"time"

	settings "github.com/krasov-rf/infobot/pkg/settings/infobot"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type InfoBotDb struct {
	*sqlx.DB

	config *settings.Config

	errChan chan error
}

func New(c *settings.Config) (*InfoBotDb, error) {
	dsn := c.DB.UrlPostgres()

	conn, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(50) // по умолчанию стоит анлимит
	conn.SetMaxIdleConns(50)
	conn.SetConnMaxLifetime(200 * time.Millisecond)
	conn.SetConnMaxIdleTime(20 * time.Second)

	DB := &InfoBotDb{
		config:  c,
		DB:      conn,
		errChan: make(chan error, 10),
	}

	return DB, nil
}

func (db *InfoBotDb) Close(ctx context.Context) error {
	return db.DB.Close()
}
