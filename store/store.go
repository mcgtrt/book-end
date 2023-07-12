package store

import "context"

const (
	DBURI      = "mongodb://localhost:27017"
	DBNAME     = "book-end"
	TestDBNAME = "book-end-test"
)

type Dropper interface {
	Drop(context.Context) error
}
