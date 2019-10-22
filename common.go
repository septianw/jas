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
package main

import (
	// "errors"
	"log"

	// "plugin"

	"strings"

	// "encoding/gob"
	// "errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"bufio"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/otiai10/copy"
	pak "github.com/septianw/jas/common"

	ty "github.com/septianw/jas/types"
	"github.com/spf13/viper"
	// "gopkg.in/gocraft/dbr.v2"
)

type Module struct {
	Bootstrap func()
	Router    func(*gin.Engine)
}

func GetConfig(key string) interface{} {
	return viper.Get(key)
}

func GetAllConfig() map[string]interface{} {
	return viper.AllSettings()
}

func GetCWD() string {
	Basepath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Printf("Fail to get current working directory : %+v\n", err)
	}

	return Basepath
}

func LoadDatabase() ty.Database {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rt := pak.ReadRuntime()
	db := pak.LoadDatabase(filepath.Join(rt.Libloc, "database.so"), rt.Dbconf)

	return db
}

func CopySchema(fpath string) error {
	rt := pak.ReadRuntime()
	return copy.Copy(fpath, rt.Migloc)
}

func CopyAllSchema(fpath string) error {
	return filepath.Walk(fpath, func(base string, finfo os.FileInfo, perr error) error {
		var reval error = nil

		if (perr == nil) && (strings.Compare("schema", path.Base(base)) == 0) {
			reval = CopySchema(base)
		} else {
			reval = perr
		}

		return reval
	})
}

// TODO: Dalam schema harus ada prefix nama module:
// - ketika module dimuat, salin semua schema yang ditemukan ke Migloc.
// - ketika module dilepas, semua schema yang sebelumnya tersalin, dihapus.
// - ketika module dimuat, insert record ke database dengan status loaded.
// FIXME: still satisfy unit testing.
func MountModule(path string) (*Module, error) {
	var mod Module
	stat, ferr := os.Stat(path)
	// os.IsNotExist(ferr)
	if ferr == nil {
		log.Printf("Exist yet: %+v", stat)
		lib := pak.LoadSo(path)
		bootsym, err := lib.Lookup("Bootstrap")
		pak.ErrHandler(err)

		routersym, err := lib.Lookup("Router")
		pak.ErrHandler(err)

		mod.Bootstrap = bootsym.(func())
		mod.Router = routersym.(func(*gin.Engine))
		return &mod, nil
	}

	return nil, ferr
}

func SetModuleMounted(moduleName string) (status bool, err error) {
	status = false
	meta, err := GetModuleMetadata(moduleName)
	if err != nil {
		log.Printf("Failed to get Module Metadata: %s", err.Error())
		return
	}

	meta.Status = "mounted"

	err = SetModuleMetadata(meta)
	if err != nil {
		log.Printf("Failed to set Module Metadata: %s", err.Error())
		return
	}

	status = true
	return
}
func SetModuleUnmounted(moduleName string) (status bool, err error) {
	status = false
	meta, err := GetModuleMetadata(moduleName)
	if err != nil {
		log.Printf("Failed to get Module Metadata: %s", err.Error())
		return
	}

	meta.Status = "loaded"

	err = SetModuleMetadata(meta)
	if err != nil {
		log.Printf("Failed to set Module Metadata: %s", err.Error())
		return
	}

	status = true
	return
}

// FIXME: belum tuntas, still satisfy unit testing.
// Load module ini akan memasukkan info module (module.toml) dalam database
// yang nantinya akan ditampilkan dalam UI.
func LoadModule(modulePath string) {
	var count int

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigName("module")
	v.AddConfigPath(modulePath)

	err := v.ReadInConfig()
	pak.ErrHandler(err)

	rt := pak.ReadRuntime()
	// cdb := LoadDatabase()

	// db, err := cdb.OpenDb(rt.Dbconf)
	// pak.ErrHandler(err)
	// log.Printf("ini db yang dimuat dari function LoadDatabase: %+v", db)

	// println("checkthis")
	log.Printf("module %s config : %+v", modulePath, v.AllSettings())
	log.Printf("sofile: %s", filepath.Join(modulePath, v.GetString("sofile")))
	qcheck := fmt.Sprintf("select count(*) from modules where name = '%s'", v.GetString("name"))
	q := fmt.Sprintf("insert into modules (name, version, status, sopath) values ('%s', '%s', '%s', '%s')",
		v.GetString("name"),
		v.GetString("version"),
		"loaded",
		filepath.Join(modulePath, v.GetString("sofile")),
	)
	dbase := LoadDatabase()
	db, err := dbase.OpenDb(rt.Dbconf)
	defer db.Close()
	pak.ErrHandler(err)

	rows, err := db.Query(qcheck)
	pak.ErrHandler(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&count)
		pak.ErrHandler(err)
	}
	log.Printf("\nquery: %+v\n", qcheck)
	log.Printf("\nres count: %+v\n", count)

	if count == 0 {
		res, err := db.Exec(q)
		pak.ErrHandler(err)
		raff, err := res.RowsAffected()
		pak.ErrHandler(err)
		log.Printf("\nres count: %+v\n", raff)
	}
}

// FIXME: Affected by this issue https://github.com/septianw/jas/issues/10
func IsModuleEnabled(moduleName string) bool {
	var modules []string
	var enabled bool = false
	var cwd string

	if wd, found := os.LookupEnv("cwd"); found {
		cwd = wd
	} else {
		cwd = GetCWD()
	}
	modenfile, err := os.Open(filepath.Join(cwd, "module-enabled"))
	pak.ErrHandler(err)

	bscanner := bufio.NewScanner(modenfile)
	bscanner.Split(bufio.ScanLines)

	for bscanner.Scan() {
		modules = append(modules, bscanner.Text())
	}

	modenfile.Close()

	foundIdx := sort.SearchStrings(modules, moduleName)
	if foundIdx != len(modules) {
		if strings.Compare(modules[foundIdx], moduleName) == 0 {
			enabled = true
		}
	}

	return enabled
}

// FIXME: still satisfy unit testing.
// TODO: fokus ke mount module
// Cari module di daftar enabled, kalau ada aktifkan, kalau tidak ada lewati.
func lmod(modtype, moduleName string) (*Module, error) {
	var mod *Module
	var err error = nil
	var modpath = filepath.Join(Modloc, modtype, moduleName)
	LoadModule(modpath)

	// TODO:
	// - cari module di daftar enabled.
	// - kalau ada aktifkan.
	if strings.Compare(modtype, "core") == 0 {
		meta, err := GetModuleMetadata(moduleName)
		pak.ErrHandler(err)
		// log.Printf("\nMeta: %+v\n", meta)
		mod, err = MountModule(meta.Sopath)
	} else {
		if IsModuleEnabled(moduleName) {
			meta, err := GetModuleMetadata(moduleName)
			pak.ErrHandler(err)
			// log.Printf("\nMeta: %+v\n", meta)
			mod, err = MountModule(meta.Sopath)
		}
	}

	return mod, err
}

func LoadCoreModule(moduleName string) (*Module, error) {
	return lmod("core", moduleName)
}

func LoadContribModule(moduleName string) (*Module, error) {
	return lmod("contrib", moduleName)
}

// Function untuk initiate module. Function ini bergantung pada Modloc
func InitiateModules(modules []os.FileInfo, moduleType string) {
	for _, module := range modules {
		// if err := CopyAllSchema(coreModule.Name()); err != nil {
		// 	pak.ErrHandler(err)
		// }
		if module.IsDir() {
			// LoadCoreModule(coreModule.Name())

			err := CopyAllSchema(filepath.Join(Modloc, moduleType, module.Name()))
			pak.ErrHandler(err)

			if m, err := lmod(moduleType, module.Name()); (err == nil) && (m != nil) {
				m.Bootstrap()
				m.Router(Routers)
			} else {
				pak.ErrHandler(err)
			}
		}
	}
}

// func LoadMiddlewares() (handlers []gin.HandlerFunc) {

// }
