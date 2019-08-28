package types

import (
	"database/sql"
)

type Db string

type Dbconf struct {
	Type     string
	Host     string
	Port     uint16
	User     string
	Pass     string
	Database string
}

type Runtime struct {
	AppName        string
	Version        string
	BuildId        string
	Stage          string
	ConfigLocation string
	Dbconf         Dbconf
	Modloc         string
	Libloc         string
}

type Database interface {
	PingDb(Dbconf) (bool, error)
	OpenDb(Dbconf) (*sql.DB, error)
	Migrate(string, Dbconf)
	SetupDb(Dbconf) bool
}
