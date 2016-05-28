package main

import (
	"log"

	"github.com/jackc/pgx"
)

// Queryer is an interface for Query
type Queryer interface {
	Query(query string, args ...interface{}) (*pgx.Rows, error)
	QueryRow(query string, args ...interface{}) *pgx.Row
}

// Execer is an interface for Exec
type Execer interface {
	Exec(query string, args ...interface{}) (pgx.CommandTag, error)
}

// TxStarter is an interface to deal with transaction
type TxStarter interface {
	Begin() (*pgx.Tx, error)
}

// TxController is an interface to deal with transaction
type TxController interface {
	Commit() error
	Rollback() error
}

// Ext is a union interface which can bind, query, and exec
type Ext interface {
	Queryer
	Execer
	TxStarter
}

// Txer is a interface for Tx
type Txer interface {
	Queryer
	Execer
	TxController
}

func insertJob(ext Txer, name string, count int) error {
	if _, err := ext.Exec(`INSERT INTO test (name, count) VALUES ($1, $2)`, name, count); err != nil {
		return err
	}
	log.Println("insertJob ok")
	return nil
}

func insertWithTxAndConJob(ext Txer, con Ext, name string, count int) error {
	if _, err := ext.Exec(`INSERT INTO test (name, count) VALUES ($1, $2)`, name, count); err != nil {
		return err
	}
	if _, err := con.Exec(`INSERT INTO test (name, count) VALUES ($1, $2)`, name+"conn", count); err != nil {
		return err
	}
	return nil
}

func updateJob(ext Txer) error {
	return nil
}

func selectJob(ext Ext) error {
	return nil
}
