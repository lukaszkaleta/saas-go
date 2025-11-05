package pg

import (
	"context"
	"fmt"
	"strings"

	pgdb "github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgState struct {
	Db          *pgdb.PgDb
	TableEntity pgdb.TableEntity
	StateColumn string
}

func NewPgState(db *pgdb.PgDb, te pgdb.TableEntity, stateColumn string) universal.State {
	return &PgState{Db: db, TableEntity: te, StateColumn: stateColumn}
}

// Name queries the current state value from the database.
// If the DB is unavailable or the query fails, it returns an empty string.
func (p *PgState) Name() string {
	if p == nil || p.Db == nil || p.Db.Pool == nil {
		return ""
	}
	col := p.StateColumn
	if col == "" {
		col = "status"
	}
	query := fmt.Sprintf("select %s from %s where id = $1", col, p.TableEntity.Name)
	var val string
	if err := p.Db.Pool.QueryRow(context.Background(), query, p.TableEntity.Id).Scan(&val); err != nil {
		return ""
	}
	return val
}

// Change updates the state column in the backing table.
func (p *PgState) Change(newState string) error {
	col := p.StateColumn
	if col == "" {
		col = "status"
	}
	query := fmt.Sprintf("update %s set %s = $1 where id = $2", p.TableEntity.Name, col)
	_, err := p.Db.Pool.Exec(context.Background(), query, newState, p.TableEntity.Id)
	if err != nil {
		return err
	}
	return nil
}

type PgTimestampState struct {
	Db          *pgdb.PgDb
	TableEntity pgdb.TableEntity
	// Order matters
	StateColumns []string
}

func NewPgTimestampState(db *pgdb.PgDb, te pgdb.TableEntity, stateColumns []string) universal.State {
	return &PgTimestampState{Db: db, TableEntity: te, StateColumns: stateColumns}
}

// Name queries the current state value from the database.
// If the DB is unavailable or the query fails, it returns an empty string.
func (p *PgTimestampState) Name() string {
	if p == nil || p.Db == nil || p.Db.Pool == nil {
		return ""
	}
	if len(p.StateColumns) == 0 {
		return ""
	}

	columns := []string{}
	for _, u := range p.StateColumns {
		columns = append(columns, u)
	}
	queryColumn := ":" + strings.Join(columns, ", ")

	query := fmt.Sprintf("select %s from %s where id = $1", queryColumn, p.TableEntity.Name)

	row := p.Db.Pool.QueryRow(context.Background(), query, p.TableEntity.Id)

	timestamps := make([]interface{}, len(p.StateColumns))
	for ts := range timestamps {
		timestamps[ts] = new(interface{})
	}
	// Scan the selected timestamp columns
	if err := row.Scan(timestamps...); err != nil {
		return ""
	}
	// Iterate over timestamps to find the last non-null value
	lastIdx := -1
	for i := range timestamps {
		if ptr, ok := timestamps[i].(*interface{}); ok && ptr != nil && *ptr != nil {
			lastIdx = i
		}
	}
	if lastIdx == -1 {
		return ""
	}
	// Derive state name from column name, e.g., "state_active" -> "active"
	colName := p.StateColumns[lastIdx]
	if strings.HasPrefix(colName, "state_") {
		return strings.TrimPrefix(colName, "state_")
	}
	return colName
}

// Change updates the state column in the backing table.
func (p *PgTimestampState) Change(newState string) error {

	col := "state_" + newState
	query := fmt.Sprintf("update %s set %s = $1 where id = $2", p.TableEntity.Name, col)
	_, err := p.Db.Pool.Exec(context.Background(), query, newState, p.TableEntity.Id)
	if err != nil {
		return err
	}
	return nil
}
