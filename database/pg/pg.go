package pg

import (
	"context"
	"fmt"
	"io"
	"io/fs"
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

func (db *PgDb) ExecuteFileFromFs(fs fs.FS, path string) error {
	open, err := fs.Open(path)
	if err != nil {
		return err
	}
	sqlStatements, err := io.ReadAll(open)
	if err != nil {
		return err
	}
	sqlArray := strings.Split(string(sqlStatements), ";")
	return db.ExecuteSqls(sqlArray)
}

func NewPg() *PgDb {

	dbUrlKey := "DATABASE_URL"
	databaseUrl := os.Getenv(dbUrlKey)
	if databaseUrl == "" {
		fmt.Fprintf(os.Stderr, "Database url is not configered, Please provide environment variable: %s\n", dbUrlKey)
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

func (db *PgDb) TableEntity(name string, id int64) TableEntity {
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
