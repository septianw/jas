package main

import (
	"fmt"

	pak "github.com/septianw/jas/common"
	ty "github.com/septianw/jas/types"

	"log"
)

func GetModuleMetadata(moduleName string) (meta ty.ModuleMetadata, err error) {
	rt = pak.ReadRuntime()
	sdb := LoadDatabase()
	db, err := sdb.OpenDb(rt.Dbconf)
	pak.ErrHandler(err)
	defer db.Close()

	q := fmt.Sprintf("select name, version, status, sopath from modules where name = '%s'",
		moduleName)
	rows, err := db.Query(q)
	pak.ErrHandler(err)

	for rows.Next() {
		err := rows.Scan(&meta.Name, &meta.Version, &meta.Status, &meta.Sopath)
		pak.ErrHandler(err)
	}

	log.Printf("ModuleMeta: %+v\n", meta)
	return
}

func SetModuleMetadata(meta ty.ModuleMetadata) (err error) {
	rt = pak.ReadRuntime()
	sdb := LoadDatabase()
	db, err := sdb.OpenDb(rt.Dbconf)
	defer db.Close()
	if err != nil {
		return
	}

	qs := `UPDATE TABLE modules SET version = '%s', status = '%s', sopath = '%s'`
	q := fmt.Sprintf(qs, meta.Version, meta.Status, meta.Sopath)

	result, err := db.Exec(q)
	if err != nil {
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return
	}

	log.Printf("%d Module updated.", affected)
	log.Printf("Module %s updated", meta.Name)

	return
}
