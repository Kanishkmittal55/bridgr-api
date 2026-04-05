package auth

type Access string

const (
	AccessAdmin Access = "ADMIN"
	AccessRead  Access = "READ"
	AccessWrite Access = "WRITE"
)
