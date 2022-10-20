package postgres

import (
	"database/sql"
	"fmt"

	"github.com/AnNosov/simple_user_api/config"
	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

func New(cfg *config.Postgres) (*Postgres, error) {

	connStr := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBname, cfg.SSLmode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	pg := &Postgres{
		DB: db,
	}

	return pg, nil

}

func (p *Postgres) Close() {
	p.DB.Close()
}
