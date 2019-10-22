package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"flag"
	"path/filepath"

	// "errors"

	"github.com/otiai10/copy"

	"github.com/septianw/jas/common"
	// "github.com/septianw/jas/types"
)

const VERSION = "0.1"

func cleanupMig(migloc string) (err error) {
	return filepath.Walk(migloc, func(base string, finfo os.FileInfo, err error) error {
		if strings.Contains(base, ".sql") {
			return os.Remove(base)
		}
		return nil
	})
}

func main() {
	// v := flag.String("v", VERSION, "Print version number")
	flag.Parse()
	if (len(flag.Args()) != 0) && (flag.Args()[0] == "version") {
		println(VERSION)
		os.Exit(0)
	}
	// println(*v)
	// if v != nil {
	// 	println(*v)
	// 	os.Exit(0)
	// }

	rt := common.ReadRuntime()
	migloc := strings.Split(rt.Migloc, "/")
	migloc[len(migloc)-1] = "schema"
	schemaloc := strings.Join(migloc, "/")

	// fmt.Printf("%+v\n", strings.Join(migloc, "/"))
	fmt.Printf("schemaloc: %+v\n", schemaloc)
	fmt.Printf("ModLoc: %+v\n", rt.Modloc)
	// filepath.w

	filepath.Walk(schemaloc, func(base string, finfo os.FileInfo, err error) error {
		if base != schemaloc {
			dst := strings.Join(append(strings.Split(rt.Migloc, "/"), finfo.Name()), "/")
			/*
				fmt.Printf("dst: %+v\n", dst)
				fmt.Printf("migloc: %+v\n", rt.Migloc)
				fmt.Printf("base: %+v\n", base)
				fmt.Printf("finfo.name(): %+v\n", finfo.Name())
				fmt.Printf("err: %+v\n", err)
			*/
			return copy.Copy(base, dst)
		}
		return nil
	})
	filepath.Walk(rt.Modloc, func(base string, finfo os.FileInfo, err error) error {
		if strings.Contains(base, ".sql") {
			dst := strings.Join(append(strings.Split(rt.Migloc, "/"), finfo.Name()), "/")
			/*
				fmt.Printf("dst: %+v\n", dst)
				fmt.Printf("base: %+v\n", base)
				fmt.Printf("finfo.name(): %+v\n", finfo.Name())
				fmt.Printf("err: %+v\n", err)
			*/
			return copy.Copy(base, dst)
		}
		return nil
	})

	// Check if database library present
	if _, err := os.Stat(filepath.Join(rt.Libloc, "database.so")); err == nil {
		db := common.LoadDatabase(filepath.Join(rt.Libloc, "database.so"), rt.Dbconf)

		common.TryCatchBlock{
			Try: func() {
				fmt.Println(" Testing database config.")
				succeed, errPing := db.PingDb(rt.Dbconf)
				if !succeed {
					log.Fatalln(errPing)
					os.Exit(3)
				}
				log.Printf("Ping database succeed: %+v\n", succeed)
			},
			Catch: func(e common.Exception) {
				log.Fatalf("Error raised while running PingDb: %+v", e)
				os.Exit(3)
			},
			Finally: func() {
				fmt.Println(" Database config test success")
			},
		}.Do()

		common.TryCatchBlock{
			Try: func() {
				fmt.Println(" Migrating database structure")
				if !db.Migrate(rt.Migloc, rt.Dbconf) {
					// fmt.Println("Database migration success.")
					fmt.Println("Database structure migration failed.")
					os.Exit(3)
				}

			},
			Catch: func(e common.Exception) {
				log.Fatalf("Error raised while running SetupDb: %+v", e)
				os.Exit(3)
			},
			Finally: func() {
				fmt.Println(" Database structure migration success")
				// cleanupMig(rt.Migloc)
			},
		}.Do()
	} else {
		log.Println(" [Warning] Database library not found, skip checking database.")
		i, e := os.Stat(filepath.Join(rt.Libloc, "database.so"))
		log.Printf("\n%+v | %+v | %+v\n", os.IsNotExist(e), i, e)
		log.Println(filepath.Join(rt.Libloc, "database.so"))
	}
}
