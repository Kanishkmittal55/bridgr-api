package repository

// Repo implements Bridgr skill-gap persistence (sqlc-backed).
type Repo struct{}

// New returns a Bridgr repository with no mutable state.
func New() *Repo {
	return &Repo{}
}
