/*
   Copyright 2019 Septian Wibisono

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
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
	MiddlewareLoc  string
	Migloc         string
}

type Database interface {
	PingDb(Dbconf) (bool, error)
	OpenDb(Dbconf) (*sql.DB, error)
	Migrate(string, Dbconf) bool
}

type ModuleMetadata struct {
	Name    string
	Version string
	Status  string
	Sopath  string
}
