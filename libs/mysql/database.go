package main

import (
	"database/sql"
	// "errors"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/septianw/jas/types"
)

type database string

// func New() database {
// 	return database{
// 		PingDb: func(d Dbconf) (bool, error) {
// 			err := errors.New("fail")
// 			return true, err
// 		},
// 		SetupDb: func(d Dbconf) bool {
// 			return true
// 		},
// 		Migrate: func(s string, d Dbconf) {
// 			return
// 		},
// 		OpenDb: func(d Dbconf) (*sql.DB, error) {
// 			return
// 		},
// 	}
// }

func (db database) PingDb(d types.Dbconf) (bool, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.User, d.Pass, d.Host, d.Port, d.Database)

	dbi, err := sql.Open("mysql", dsn)
	if err != nil {
		return false, err
	}

	err = dbi.Ping()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (db database) OpenDb(d types.Dbconf) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s/%s@%s:%d/%s",
		d.User, d.Pass, d.Host, d.Port, d.Database)

	database, err := sql.Open("goracle", dsn)

	fmt.Printf("\n%+v   %+v\n", d, database)

	return database, err
}

func (db database) Migrate(location string, d types.Dbconf) {
	fmt.Printf("\n%+v  %+v\n", location, d)
	return
}

func (db database) SetupDb(d types.Dbconf) bool {
	fmt.Printf("\n%+v\n", d)
	return true
}

var Database database
