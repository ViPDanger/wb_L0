package postgres

import (
	"context"
	"log"
	"time"

	"github.com/ViPDanger/L0/Server/Internal/config"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, conf config.CFG) (pool *pgx.Conn, err error) {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	url := "postgres://" + conf.PG_user + ":" + conf.PG_password + "@" + conf.PG_host + ":" + conf.PG_port + "/" + conf.PG_bdname
	attempts := conf.Con_Attempts
	for attempts > 0 {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		pool, err = pgx.Connect(ctx, url)
		if err != nil {
			log.Println("Connecting to Postgres: Server didn't respound")
			time.Sleep(5 * time.Second)
		}
		attempts--
	}
	if err != nil {
		log.Fatalln("Connecting to Postgres: Number of attempts exceeded.")
	}
	log.Println("Connecting to Postgres: Successfully ")
	return pool, nil
}
