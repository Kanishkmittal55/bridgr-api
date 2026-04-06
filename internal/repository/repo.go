package repository

import "github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"

// Querier is the sqlc-generated interface (skill-gap, job discovery, feed pipeline, supported boards). Satisfied by *sqlc.Queries.
type Querier = sqlc.Querier

// Repo implements Bridgr persistence (sqlc-backed): skill-gap, job discovery, feed pipeline, supported job boards.
type Repo struct{}

// New returns a Bridgr repository with no mutable state.
func New() *Repo {
	return &Repo{}
}
