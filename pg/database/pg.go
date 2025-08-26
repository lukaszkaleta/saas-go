package database

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgDb struct {
	Pool *pgxpool.Pool
}

func (db *PgDb) ExecuteSqls(sqls []string) error {
	for _, sql := range sqls {
		err := db.ExecuteSql(sql)
		ifPanic(err)
	}
	return nil
}

func (db *PgDb) ExecuteSql(sql string) error {
	_, err := db.Pool.Exec(context.Background(), sql)
	return err
}

func NewPg() *PgDb {

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		fmt.Fprintf(os.Stderr, "Database url is not configered, Please provide environment variable: %v\n")
		os.Exit(1)
	}
	dbpool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create DATABASE connection pool: %v\n", err)
		os.Exit(1)
	}

	var version string
	err = dbpool.QueryRow(context.Background(), "select version()").Scan(&version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	return &PgDb{Pool: dbpool}
}

func (db *PgDb) tableEntity(name string, id int64) TableEntity {
	return TableEntity{Name: name, Id: id}
}

func ExecuteFromFile(path string) {
	sqlStatements, err := os.ReadFile(path)
	ifPanic(err)
	sqlArray := strings.Split(string(sqlStatements), ";")
	ifPanic(NewPg().ExecuteSqls(sqlArray))
}

func ifPanic(e error) {
	if e != nil {
		panic(e)
	}
}
