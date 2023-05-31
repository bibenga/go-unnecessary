package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var db *sql.DB

func connect() (*sql.DB, error) {
	dsn := "host=db user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=UTC"
	db2, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db2.SetMaxIdleConns(0)
	db2.SetMaxOpenConns(1)
	return db2, nil
}

func trx(t *testing.T, readOnly bool) *sql.Tx {
	t.Helper()
	assert := require.New(t)

	assert.NotNil(db)
	assert.NoError(db.Ping())

	t.Log("tx - Begin")
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  readOnly,
	})
	assert.NoError(err)

	t.Cleanup(func() {
		t.Helper()
		t.Log("tx - Rollback")
		err := tx.Rollback()
		assert.NoError(err)
	})
	return tx
}

func trx2(t *testing.T, readOnly bool) *sql.Tx {
	t.Helper()
	assert := require.New(t)

	t.Log("db - Open")
	db2, err := connect()
	assert.NoError(err)
	assert.NotNil(db2)
	assert.NoError(db2.Ping())

	t.Log("tx - Begin")
	tx, err := db2.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  readOnly,
	})
	assert.NoError(err)

	// выполняется в обратном порядке
	t.Cleanup(func() {
		t.Helper()
		t.Log("db - Close")
		err = db2.Close()
		assert.NoError(err)
	})
	t.Cleanup(func() {
		t.Helper()
		t.Log("tx - Rollback")
		err := tx.Rollback()
		assert.NoError(err)
	})
	return tx
}

func TestMain(m *testing.M) {
	log.Printf("> TestMain")
	defer log.Printf("< TestMain")

	var err error
	db, err = connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	os.Exit(m.Run())
}

func TestConnected(t *testing.T) {
	assert := require.New(t)
	// if db == nil {
	// 	t.Fatal("db is nil")
	// }
	assert.NotNil(db)

	// if err := db.Ping(); err != nil {
	// 	t.Fatal(err)
	// }
	err := db.Ping()
	assert.NoError(err)

	// if _, err := db.Exec("select 1"); err != nil {
	// 	t.Fatal(err)
	// }
	row := db.QueryRow("select 1")
	assert.NoError(row.Err())
	assert.NotNil(row)
	var cnt int = -1
	assert.NoError(row.Scan(&cnt))
	assert.Equal(cnt, 1)
}

func TestTrx(t *testing.T) {
	assert := require.New(t)
	tx := trx2(t, false)
	assert.NotNil(tx)

	row := tx.QueryRow("select 1,2")
	assert.NoError(row.Err())
	assert.NotNil(row)
	var v1, v2 int
	assert.NoError(row.Scan(&v1, &v2))
	assert.Equal(v1, 1)
	assert.Equal(v2, 2)
}

func TestTag(t *testing.T) {
	assert := require.New(t)

	log.Printf("> TestTag")
	defer log.Printf("< TestTag")

	tx := trx(t, false)

	tag, err := getTag(context.Background(), tx, "olala")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if tag != nil {
	// 	t.Fatalf("found tag - %v", tag)
	// }
	assert.NoError(err)
	assert.Nil(tag)

	tag, err = insertTag(context.Background(), tx, "olala")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if tag == nil {
	// 	t.Fatal("tag is nil")
	// } else {
	// 	if tag.ID <= 0 {
	// 		t.Fatalf("bad tag ID is invalid - %v", tag)
	// 	}
	// 	if tag.Name != "OLALA" {
	// 		t.Fatalf("bad tag name - %v", tag)
	// 	}
	// }
	assert.NoError(err)
	assert.NotNil(tag)
	assert.Greater(tag.ID, int64(0))
	assert.Equal(tag.Name, "OLALA")

	affected, err := deleteTag(context.Background(), tx, "olala")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if affected != 1 {
	// 	t.Fatalf("invalid affected - %v", affected)
	// }
	assert.NoError(err)
	assert.Equal(affected, int64(1))
}
