package bridgr

import (
	"log"

	"github.com/hassleskip/bridgr-api/internal/api/deps"
)

type server struct {
	deps *deps.Deps
}

func NewServer(d *deps.Deps) *server {
	if d == nil {
		log.Panic("bridgr.NewServer: deps must not be nil")
	}
	if d.ResponseWriter == nil {
		log.Panic("bridgr.NewServer: deps.ResponseWriter must not be nil")
	}
	return &server{deps: d}
}
