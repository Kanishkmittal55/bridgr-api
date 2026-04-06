package main

import (
	"log"

	"github.com/Kanishkmittal55/bridgr-api/internal/config"
	"github.com/Kanishkmittal55/bridgr-api/internal/rdbms"
	"github.com/Kanishkmittal55/bridgr-api/internal/seeds"
)

func main() {
	cfg := config.Load()
	db, err := rdbms.NewConn(rdbms.ConnStr(cfg))
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	seeds.Execute(db)
}
