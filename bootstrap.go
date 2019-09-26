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
	"fmt"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"

	// "errors"
	"io/ioutil"

	// "reflect"
	"strings"

	"path/filepath"

	pak "github.com/septianw/jas/common"
	ty "github.com/septianw/jas/types"

	// ora "github.com/septianw/shiny-telegram/experiment01/sharedpak/dboracle"
	"github.com/spf13/viper"
)

const BOOTSTRAP_LEVEL_0 = 0
const BOOTSTRAP_LEVEL_1 = 1
const BOOTSTRAP_LEVEL_2 = 2
const BOOTSTRAP_LEVEL_3 = 3

var Spin = spinner.New(spinner.CharSets[24], 100*time.Millisecond)
var ListenAddr, Dsn string
var rt ty.Runtime

// var Config

// NOTE: Dari setiap module ada semacam hook yang dapat dipanggil pada bootstrap level berapa.

// check integrity (rely on system, we can't check ourself id)
// check requirement
//   paths
//   config
//   libraries
func RunBootLevel0() {
	// var files []string

	// fmt.Println()
	Spin.Start()
	Spin.Suffix = "  Check files existence:"

	rt.AppName = APPNAME
	rt.BuildId = BUILDID
	rt.Version = VERSION

	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", APPNAME))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", APPNAME))
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalln("No config file found.")
			// log.Printf("err: %+v\n", err)
			os.Exit(2)
			// Config file not found; ignore error if desired
		} else {
			log.Fatalln(err)
			os.Exit(2)
			// Config file was found but another error was produced
		}
	}

	// fmt.Printf("\n%+v\n", viper.Get("schema"))

	// fmt.Printf("Loaded configuration file: %s\n", viper.ConfigFileUsed())
	rt.ConfigLocation = viper.ConfigFileUsed()

	ListenAddr = fmt.Sprintf("%s:%d", viper.GetString("server.bind"), viper.GetInt("server.port"))
	// fmt.Printf("Listening at %s\n", ListenAddr)

	switch os.Getenv("STAGE") {
	case "production":
		STAGE = "production"
	case "development":
		fallthrough
	case "testing":
		fallthrough
	default:
		STAGE = "development"
	}
	rt.Stage = STAGE

	// Load main library
	Libloc = viper.GetString("libraryLocation")

	// // if modloc is empty use default location Current Working Directory
	if strings.Compare(Libloc, "") == 0 {
		if strings.Compare(LIBRARY_LOCATION, "") != 0 {
			Libloc = LIBRARY_LOCATION
		} else {
			cwd, err := os.Getwd()
			pak.ErrHandler(err)
			Libloc = filepath.Join(cwd, "libs")
		}
	}
	rt.Libloc = Libloc
	pak.WriteRuntime(rt)

	// time.Sleep(10 * time.Second)
	Spin.Stop()
}

// basic connectivity
//   db
//   cache
// basic table structure
//   check schema structure
//   schema exist
func RunBootLevel1() {
	var dbconf ty.Dbconf
	var db ty.Database

	RunBootLevel0()
	var rt = pak.ReadRuntime()
	Spin.Start()
	Spin.Suffix = "  This is booting level 1"

	// basic connectivity
	// TODO: Semua yang ada comment, itu yang sebelumnya berjalan dengan baik.
	//       Sampai negara api menyerang.
	// d := GetConfig("database").(map[string]interface{})
	// fmt.Printf("|%+v|", reflect.TypeOf(d["hostname"]))
	// d := viper.Get("database").(map[string]interface{})

	Migloc = viper.GetString("migrationLocation")
	if strings.Compare(Migloc, "") == 0 {
		if strings.Compare(LIBRARY_LOCATION, "") != 0 {
			Migloc = LIBRARY_LOCATION
		} else {
			cwd, err := os.Getwd()
			pak.ErrHandler(err)
			Migloc = filepath.Join(cwd, "migrations")
		}
	}
	rt.Migloc = Migloc

	// Check if database library present
	if _, err := os.Stat(filepath.Join(Libloc, "database.so")); err == nil {
		dbconf.Host = viper.GetString("database.hostname") // d["hostname"].(string)
		dbconf.Type = viper.GetString("database.type")     // d["type"].(string)
		// convert dari map viper ke int64 dan convert lagi ke uint16
		// karena int64 terlalu besar untuk menyimpan port yang isinya maksimum hanya 65535
		dbconf.Port = uint16(viper.GetInt64("database.port"))  // uint16(d["port"].(int64))
		dbconf.User = viper.GetString("database.username")     // d["username"].(string)
		dbconf.Pass = viper.GetString("database.password")     // d["password"].(string)
		dbconf.Database = viper.GetString("database.database") // d["database"].(string)

		db = pak.LoadDatabase(filepath.Join(Libloc, "database.so"), dbconf)

		pak.TryCatchBlock{
			Try: func() {
				Spin.Suffix = " Testing database config"
				succeed, errPing := db.PingDb(dbconf)
				if !succeed {
					log.Fatalln(errPing)
					os.Exit(3)
				}
				log.Printf("Ping database succeed: %+v\n", succeed)
			},
			Catch: func(e pak.Exception) {
				log.Fatalf("Error raised while running PingDb: %+v", e)
				os.Exit(3)
			},
			Finally: func() {
				Spin.Suffix = " Database config test success"
			},
		}.Do()

		pak.TryCatchBlock{
			Try: func() {
				Spin.Suffix = " Migrating database structure"
				if !db.Migrate(rt.Migloc, dbconf) {
					// fmt.Println("Database migration success.")
					fmt.Println("Database structure migration failed.")
					os.Exit(3)
				}

			},
			Catch: func(e pak.Exception) {
				log.Fatalf("Error raised while running SetupDb: %+v", e)
				os.Exit(3)
			},
			Finally: func() {
				Spin.Suffix = " Database structure migration success"
			},
		}.Do()

		rt.Dbconf = dbconf
	} else {
		log.Println(" [Warning] Database library not found, skip checking database.")
		i, e := os.Stat(filepath.Join(Libloc, "database.so"))
		log.Printf("\n%+v | %+v | %+v\n", os.IsNotExist(e), i, e)
		log.Println(filepath.Join(Libloc, "database.so"))
	}

	pak.WriteRuntime(rt)

	Spin.Stop()
}

// TODO: tambahkan config manifest pada setiap module
// TODO: load config config itu dan gunakan viper merge config untuk merge.
// TODO: format config pakai map, lalu loop config tiap module pakai range map.

// init core
//   setup
//   run
// collecting module
// setup basic module
func RunBootLevel2() {
	RunBootLevel1()
	// var modules []*Module
	var rt = pak.ReadRuntime()

	Spin.Start()
	Spin.Suffix = "  This is booting level 2"

	// Load Core module
	Modloc = viper.GetString("moduleLocation")

	// if modloc is empty use default location Current Working Directory
	if strings.Compare(Modloc, "") == 0 {
		if strings.Compare(MODULE_LOCATION, "") != 0 {
			Modloc = MODULE_LOCATION
		} else {
			cwd, err := os.Getwd()
			pak.ErrHandler(err)
			Modloc = filepath.Join(cwd, "modules")
		}
	}
	// log.Println(Modloc)

	rt.Modloc = Modloc

	Spin.Suffix = " Initiate core modules"
	coreModules, err := ioutil.ReadDir(filepath.Join(Modloc, "core"))
	pak.ErrHandler(err)

	// Setup router for the first time.
	Routers = SetupRouter()

	// TODO: ada beberapa skenario disini:
	// 1. muat semua module dalam direktori secara langsung.
	//    semua module akan terinstall dan termuat secara otomatis.
	// 2. memuat, menginstall, dan uninstall module dilakukan manual,
	//    melalui sebuah endpoint untuk melakukan operasi itu.
	for _, coreModule := range coreModules {
		// if err := CopyAllSchema(coreModule.Name()); err != nil {
		// 	pak.ErrHandler(err)
		// }
		if coreModule.IsDir() {
			// LoadCoreModule(coreModule.Name())
			if m, err := LoadCoreModule(coreModule.Name()); (err == nil) && (m != nil) {
				m.Bootstrap()
				m.Router(Routers)
			} else {
				pak.ErrHandler(err)
			}
		}
	}
	// Load Core module done

	// fmt.Printf("%+v", Modloc)

	// time.Sleep(10 * time.Second)
	pak.WriteRuntime(rt)
	Spin.Stop()
}

// init contrib
//   setup
//   run
// setup router
func RunBootLevel3() {
	RunBootLevel2()
	Spin.Start()
	Spin.Suffix = "  This is booting level 3"

	Modloc = viper.GetString("moduleLocation")

	// if modloc is empty use default location or Current Working Directory
	if strings.Compare(Modloc, "") == 0 {
		if strings.Compare(MODULE_LOCATION, "") != 0 {
			Modloc = MODULE_LOCATION
		} else {
			cwd, err := os.Getwd()
			pak.ErrHandler(err)
			Modloc = cwd + "/modules"
		}
	}

	Spin.Suffix = " Initiate contributed modules"
	contribModules, err := ioutil.ReadDir(strings.Join(
		[]string{Modloc, "contrib"}, "/"))
	pak.ErrHandler(err)

	// Setup router for the first time.
	// Routers = SetupRouter()

	for _, contribModule := range contribModules {
		if contribModule.IsDir() {
			if err := CopyAllSchema(contribModule.Name()); err != nil {
				pak.ErrHandler(err)
			}
			if m, err := LoadContribModule(contribModule.Name()); err == nil {
				m.Bootstrap()
				m.Router(Routers)
			}
		}
	}

	// time.Sleep(10 * time.Second)
	Spin.Stop()
}

func Bootstrap(level int) {
	switch level {
	case BOOTSTRAP_LEVEL_0:
		RunBootLevel0()
		break
	case BOOTSTRAP_LEVEL_1:
		RunBootLevel1()
		break
	case BOOTSTRAP_LEVEL_2:
		RunBootLevel2()
		break
	case BOOTSTRAP_LEVEL_3:
		RunBootLevel3()
		break
	}
	pak.WriteRuntime(rt)
}

func BootstrapAll() {
	Bootstrap(BOOTSTRAP_LEVEL_3)
}
