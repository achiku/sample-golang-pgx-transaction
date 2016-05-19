package main

import (
	"testing"

	"github.com/jackc/pgx"
)

var testConnConfig = pgx.ConnConfig{
	Host:     "localhost",
	User:     "pgtest",
	Database: "pgtest",
}

func testOpenConnPool(t *testing.T, maxConnections int) *pgx.ConnPool {
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig:     testConnConfig,
		MaxConnections: maxConnections,
	}
	pool, err := pgx.NewConnPool(connPoolConfig)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec("CREATE TABLE IF NOT EXISTS test (name text primary key, count integer)"); err != nil {
		t.Fatal(err)
	}
	return pool
}

func testDropAndClose(t *testing.T, pool *pgx.ConnPool) {
	if _, err := pool.Exec("drop TABLE test"); err != nil {
		t.Fatal(err)
	}
	pool.Close()
}

func TestInsertJob(t *testing.T) {
	conn := testOpenConnPool(t, 2)
	tx, err := conn.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	defer tx.Rollback()

	name := "test1"
	count := 10
	if err := insertJob(tx, name, count); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := tx.QueryRow(`select count(*) from test`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Errorf("want 1 got %d", n)
	}
	var c int
	var nm string
	if err := tx.QueryRow(`select name, count from test`).Scan(&nm, &c); err != nil {
		t.Fatal(err)
	}
	if nm != name {
		t.Errorf("want %s got %s", name, nm)
	}
	if c != count {
		t.Errorf("want %d got %d", count, c)
	}
}

func TestInsertWithTxAndConJob(t *testing.T) {
	conn := testOpenConnPool(t, 2)
	tx, err := conn.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	defer tx.Rollback()

	name := "test1"
	count := 10
	if err := insertWithTxAndConJob(tx, conn, name, count); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := tx.QueryRow(`select count(*) from test`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Errorf("want 2 got %d", n)
	}
	var c int
	var nm string
	if err := tx.QueryRow(`select name, count from test where name like '%conn'`).Scan(&nm, &c); err != nil {
		t.Fatal(err)
	}
	if nm != name {
		t.Errorf("want %s got %s", name, nm)
	}
	if c != count {
		t.Errorf("want %d got %d", count, c)
	}
}

func TestSelectJob(t *testing.T) {
}
