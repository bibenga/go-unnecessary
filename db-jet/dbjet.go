package main

// https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md

//go:generate jet -dsn=postgresql://postgres:postgres@db:5432/postgres?sslmode=disable -path=./ -ignore-tables schema_migrations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"

	jet "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	m "unnecessary/db-jet/go/public/model"
	t "unnecessary/db-jet/go/public/table"
)

const dsn string = "host=host.docker.internal port=5432 user=rds password=sqlsql dbname=go TimeZone=UTC"

func deleteTag(ctx context.Context, q qrm.DB, name string) (int64, error) {
	log.Printf("delete tag '%s'", name)

	tagDeleteStmt := t.Tag.DELETE().WHERE(
		jet.LOWER(t.Tag.Name).EQ(jet.LOWER(jet.String("TAG1"))),
	)
	log.Printf("run SQL %s", tagDeleteStmt.DebugSql())
	res, err := tagDeleteStmt.ExecContext(ctx, q)
	if err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	log.Printf("affected - %+v", affected)
	return affected, nil
}

func getTag(ctx context.Context, q qrm.DB, name string) (*m.Tag, error) {
	log.Printf("get tag '%s'", name)

	tagSelectStmt := jet.SELECT(
		t.Tag.AllColumns,
	).FROM(
		t.Tag,
	).WHERE(
		jet.LOWER(t.Tag.Name).EQ(jet.LOWER(jet.String(name))),
	)

	log.Printf("run SQL %s", tagSelectStmt.DebugSql())
	var tag m.Tag
	err := tagSelectStmt.Query(q, &tag)
	if errors.Is(err, qrm.ErrNoRows) {
		log.Printf("tag '%s' not found", name)
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		log.Printf("tag '%s' found with %v", name, tag.ID)
		return &tag, nil
	}
}

func insertTag(ctx context.Context, q qrm.DB, name string) (*m.Tag, error) {
	log.Printf("create new tag '%s'", name)

	tag := m.Tag{
		Name:       strings.ToUpper(name),
		ModifiedTs: time.Now(),
	}
	log.Printf("insert: %+v", tag)
	tagInsertStmt := t.Tag.INSERT(
		t.Tag.Name, t.Tag.ModifiedTs,
	).MODEL(
		tag,
	).RETURNING(
		t.Tag.ID,
	)
	log.Printf("run SQL %s", tagInsertStmt.DebugSql())
	var itag m.Tag
	err := tagInsertStmt.Query(q, &itag)
	if err != nil {
		return nil, err
	}
	tag.ID = itag.ID
	log.Printf("inserted: %+v", tag)
	return &tag, nil
}

func updateTag(ctx context.Context, q qrm.DB, tag *m.Tag) (int64, error) {
	log.Printf("update: %+v", *tag)
	tag.ModifiedTs = time.Now()

	tagUpdateStmt := t.Tag.UPDATE(
		t.Tag.ModifiedTs,
	).MODEL(
		tag,
	).WHERE(
		t.Tag.ID.EQ(jet.Int64(tag.ID)),
	)
	log.Printf("run SQL %s", tagUpdateStmt.DebugSql())
	res, err := tagUpdateStmt.ExecContext(ctx, q)
	if err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	log.Printf("affected - %v", affected)
	log.Printf("updated: %+v", *tag)
	return 1, nil
}

func playWithDbConJet() {
	// https://github.com/go-jet/jet
	// go get -u github.com/go-jet/jet/v2
	// go install github.com/go-jet/jet/v2/cmd/jet@latest
	log.Printf("playWithDbConJet")

	// db, err := sql.Open("postgres", dsn)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	tagSelectStmt := jet.SELECT(
		t.Tag.AllColumns,
	).FROM(
		t.Tag,
	).WHERE(
		jet.LOWER(t.Tag.Name).EQ(jet.LOWER(jet.String("TAG1"))),
	)

	var tags []m.Tag
	// log.Debug().Interface("sql", stmt.DebugSql()).Msgf("Run SQL")
	log.Printf("run SQL %s", tagSelectStmt.DebugSql())
	err = tagSelectStmt.Query(tx, &tags)
	if err != nil {
		panic(err)
	}

	var tag m.Tag
	if len(tags) == 0 {
		log.Print("Insert new tag")
		tag = m.Tag{
			Name:       "TAG1",
			ModifiedTs: time.Now(),
		}
		tagInsertStmt := t.Tag.INSERT(
			t.Tag.Name, t.Tag.ModifiedTs,
		).MODEL(
			tag,
		).RETURNING(
			t.Tag.ID,
		)
		log.Printf("run SQL %s", tagInsertStmt.DebugSql())
		err = tagInsertStmt.Query(tx, &tags)
		if err != nil {
			panic(err)
		}
		if len(tags) != 1 {
			panic(fmt.Sprintf("tags count %d!", len(tags)))
		}
		tag.ID = tags[0].ID
		log.Printf("Inserted tag: %+v", tag)

	} else if len(tags) == 1 {
		log.Printf("Tag found")
		tag = tags[0]
		log.Printf("Loaded tag: %+v", tag)

		tag.ModifiedTs = time.Now()

		tagUpdateStmt := t.Tag.UPDATE(
			t.Tag.ModifiedTs,
		).MODEL(
			tag,
		).WHERE(
			t.Tag.ID.EQ(jet.Int64(tag.ID)),
		)
		log.Printf("run SQL %s", tagUpdateStmt.DebugSql())
		err = tagUpdateStmt.Query(tx, &tags)
		if err != nil {
			panic(err)
		}
		log.Printf("Updated tag: %+v", tag)

	} else {
		panic(fmt.Sprintf("tags count %d!", len(tags)))
	}

	// commit
	if err := tx.Commit(); err != nil {
		panic("can't close transaction")
	}
}

func playWithDbConJet2() {
	// https://github.com/go-jet/jet
	// go get -u github.com/go-jet/jet/v2
	// go install github.com/go-jet/jet/v2/cmd/jet@latest
	log.Printf("playWithDbConJet")

	// db, err := sql.Open("postgres", dsn)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	log.Print("----------------------")
	_, err = deleteTag(context.Background(), tx, "tag1")
	if err != nil {
		panic(err)
	}
	log.Print("----------------------")
	_, err = getTag(context.Background(), tx, "tag1")
	if err != nil {
		panic(err)
	}
	log.Print("----------------------")
	tag, err := insertTag(context.Background(), tx, "tag1")
	if err != nil {
		panic(err)
	}
	log.Print("----------------------")
	_, err = updateTag(context.Background(), tx, tag)
	if err != nil {
		panic(err)
	}
	log.Print("----------------------")
	_, err = getTag(context.Background(), tx, "tag1")
	if err != nil {
		panic(err)
	}

}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("jet - ")

	log.Print("start")
	// playWithDbConJet()
	playWithDbConJet2()
}
