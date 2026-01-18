package pg

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgDb struct {
	Pool *pgxpool.Pool
}

func (db *PgDb) ExecuteSqls(sqls []string) error {
	for _, sql := range sqls {
		slog.Info("Executing", "SQL", sql)
		err := db.ExecuteSql(sql)
		if err != nil {
			slog.Error("Check", "SQL", sql)
			return err
		}
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

func NewPgWithUrl(databaseUrl string) *PgDb {

	dbpool, err := pgxpool.NewWithConfig(context.Background(), Config(databaseUrl))
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

func LocalPgWithName(user string, dbName string) *PgDb {
	return NewPgWithUrl(fmt.Sprintf("postgresql://%s:%s@localhost:5432/%s", user, user, dbName))
}

func NewPg() *PgDb {
	dbUrlKey := "DATABASE_URL"
	databaseUrl := os.Getenv(dbUrlKey)
	if databaseUrl == "" {
		fmt.Fprintf(os.Stderr, "Database url is not configered, Please provide environment variable: %s\n", dbUrlKey)
		os.Exit(1)
	}
	return NewPgWithUrl(databaseUrl)
}

func (db *PgDb) TableEntity(name string, id int64) TableEntity {
	return TableEntity{Name: name, Id: id}
}

func ExecuteFromFile(path string) error {
	sqlStatements, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	sqlArray := strings.Split(string(sqlStatements), ";")
	return NewPg().ExecuteSqls(sqlArray)
}

func Config(url string) *pgxpool.Config {
	const defaultMaxConns = int32(40)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	dbConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		slog.Error("Failed to create a config", "error", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		slog.Debug("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		slog.Debug("After releasing the connection pool to the database!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		slog.Debug("Closed the connection pool to the database!!")
	}

	return dbConfig
}
