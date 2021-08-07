package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

const (
	defaultDBName  = "performance_schema"
	defaultQuery56 = `select
  COUNT_STAR AS cnt,
  SUM_TIMER_WAIT/1e12 AS sum,
  MIN_TIMER_WAIT/1e12 AS min,
  AVG_TIMER_WAIT/1e12 AS avg,
  MAX_TIMER_WAIT/1e12 AS max,
  SUM_LOCK_TIME/1e12 AS sumLock,
  SUM_ROWS_SENT AS sumRows,
  ifnull((SUM_ROWS_SENT / nullif(COUNT_STAR,0)),0) AS avgRows,
  SCHEMA_NAME AS db,
  DIGEST AS digest
from events_statements_summary_by_digest
where schema_name <> ?
order by SUM_TIMER_WAIT desc`
	defaultQuery80 = `select
  COUNT_STAR AS cnt,
  SUM_TIMER_WAIT/1e12 AS sum,
  MIN_TIMER_WAIT/1e12 AS min,
  AVG_TIMER_WAIT/1e12 AS avg,
  MAX_TIMER_WAIT/1e12 AS max,
  SUM_LOCK_TIME/1e12 AS sumLock,
  SUM_ROWS_SENT AS sumRows,
  ifnull((SUM_ROWS_SENT / nullif(COUNT_STAR,0)),0) AS avgRows,
  SCHEMA_NAME AS db,
  QUERY_SAMPLE_TEXT AS query
from events_statements_summary_by_digest
where schema_name <> ?
order by SUM_TIMER_WAIT desc`
)

func getRows(ctx context.Context, db *sql.DB, dbname, query string) (*sql.Rows, error) {
	if query == "" {
		rows, err := db.QueryContext(ctx, defaultQuery80, dbname)
		// success
		if err == nil {
			return rows, err
		}
		return db.QueryContext(ctx, defaultQuery56, dbname)
	}
	return db.QueryContext(ctx, query)
}

// Output はDBからPerformanceSchemaを読み込んでTSV形式で出力
func Output(ctx context.Context, w io.Writer, dsn, dbname, query string) error {

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}
	defer db.Close()

	rows, err := getRows(ctx, db, dbname, query)
	if err != nil {
		return fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %w", err)
	}

	cw := csv.NewWriter(w)
	cw.Comma = '\t'
	defer cw.Flush()

	if err := cw.Write(cols); err != nil {
		return fmt.Errorf("failed to write columns: %w", err)
	}

	tsv := make([]string, len(cols))
	values := make([]interface{}, len(cols))
	for i := range values {
		values[i] = &tsv[i]
	}
	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			return fmt.Errorf("failed to scan rows: %w", err)
		}

		if err := cw.Write(tsv); err != nil {
			return fmt.Errorf("failed to write values: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("db error: %w", err)
	}
	if err := cw.Error(); err != nil {
		return fmt.Errorf("csv error: %w", err)
	}
	return nil
}

func main() {
	var host, sock string
	var query string

	mysqlConfig := mysql.NewConfig()

	flag.StringVar(&mysqlConfig.User, "user", "", "Username")
	flag.StringVar(&mysqlConfig.User, "u", "", "Username (shorthand)")
	flag.StringVar(&mysqlConfig.Passwd, "password", "", "Password")
	flag.StringVar(&mysqlConfig.Passwd, "p", "", "Password (shorthand)")
	flag.StringVar(&host, "host", "", "Host")
	flag.StringVar(&host, "h", "", "Host (shorthand)")
	flag.StringVar(&sock, "socket", "", "Socket path")
	flag.StringVar(&sock, "S", "", "Socket path (shorthand)")
	flag.StringVar(&query, "execute", "", "execute command")
	flag.StringVar(&query, "e", "", "execute command (shorthand)")
	flag.Parse()

	if host != "" {
		mysqlConfig.Net = "tcp"
		mysqlConfig.Addr = host
	} else if sock != "" {
		mysqlConfig.Net = "unix"
		mysqlConfig.Addr = sock
	}

	dbname := flag.Arg(0)
	if dbname == "" {
		dbname = defaultDBName
	}
	mysqlConfig.DBName = dbname

	dsn := mysqlConfig.FormatDSN()
	if err := Output(context.Background(), os.Stdout, dsn, dbname, query); err != nil {
		log.Fatal(err)
	}
}
