package database

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // Register the MySQL driver
	"github.com/pkg/errors"
)

type MysqlProvider struct {
	*sql.DB
}

func NewMysqlProvider() *MysqlProvider {
	db, err := sql.Open("mysql", "user:password@tcp(host:port)/database_name")
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	return &MysqlProvider{
		db,
	}
}

// Query used to query data from database, select
func (mysqlProvider *MysqlProvider) Query(query string, args ...any) (*sql.Rows, error) {
	conn, err := mysqlProvider.DB.Conn(context.Background())
	if err != nil {
		return nil, errors.WithMessage(err, "get connection from pool failed")
	}
	defer conn.Close()

	rows, err := conn.QueryContext(context.Background(), query, args)
	if err != nil {
		return nil, errors.WithMessage(err, "query data from database failed")
	}

	return rows, nil
}

// Exec executes sql using connection in pool
func (mysqlProvider *MysqlProvider) Exec(sql string, args ...any) (sql.Result, error) {
	conn, err := mysqlProvider.DB.Conn(context.Background())
	if err != nil {
		return nil, errors.WithMessage(err, "get connection from pool failed")
	}
	defer conn.Close()

	result, err := conn.ExecContext(context.Background(), sql, args)
	if err != nil {
		return nil, errors.WithMessage(err, "execute sql failed")
	}
	return result, nil
}
