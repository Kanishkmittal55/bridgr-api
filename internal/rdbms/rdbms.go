package rdbms

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hassleskip/bridgr-api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewConn opens a pooled Postgres connection.
func NewConn(connStr string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}
	return db, nil
}

// ConnStr builds a pgx connection string from config.
func ConnStr(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s@%s:%s/%s%s",
		usernameAndPassword(cfg.PostgresUser, cfg.PostgresPassword),
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDb,
		pgxCnxStringParams(cfg),
	)
}

func pgxCnxStringParams(cfg *config.Config) string {
	var sb strings.Builder
	if cfg.PostgresSslModeDisabled {
		sb.WriteString("?sslmode=disable")
	} else {
		sb.WriteString("?sslmode=require")
	}
	sb.WriteString(fmt.Sprintf("&pool_max_conns=%v", cfg.PostgresMaxOpenConn))
	sb.WriteString(fmt.Sprintf("&pool_min_conns=%v", cfg.PostgresMinIdleConn))
	sb.WriteString(fmt.Sprintf("&pool_max_conn_lifetime=%v", cfg.PostgresMaxConnLifetime))
	sb.WriteString(fmt.Sprintf("&pool_max_conn_idle_time=%v", cfg.PostgresMaxConnIdleTime))
	return sb.String()
}

func usernameAndPassword(username, password string) string {
	if password == "" {
		return username
	}
	return fmt.Sprintf("%s:%s", username, password)
}
