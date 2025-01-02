package database

import (
	"database/sql"
)

type Database interface {
	Close() error
	ExecSQLFile(file string) error
	Exec(query string, args ...any) error
	Query(query string, args ...any) (*sql.Rows, error)
	PrepareExec(query string) (*sql.Stmt, error)
	PrepareQuery(query string) (*sql.Stmt, error)
	Begin() error
	Commit() error
	Rollback() error
}
