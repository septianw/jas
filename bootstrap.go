package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"

	"io/ioutil"
	// "reflect"
	"strings"

	pak "github.com/septianw/jas/common"
	ora "github.com/septianw/shiny-telegram/experiment01/sharedpak/dboracle"
	"github.com/spf13/viper"
)

const BOOTSTRAP_LEVEL_0 = 0
const BOOTSTRAP_LEVEL_1 = 1
const BOOTSTRAP_LEVEL_2 = 2
const BOOTSTRAP_LEVEL_3 = 3

var Spin = spinner.New(spinner.CharSets[24], 100*time.Millisecond)
var ListenAddr, Dsn string
var rt pak.Runtime

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
			Libloc = cwd + "/lib"
		}
	}
	rt.Libloc = Libloc

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
	var dbconf pak.Dbconf

	RunBootLevel0()
	Spin.Start()
	Spin.Suffix = "  This is booting level 1"
	fmt.Printf("server: %+v", GetConfig("server"))

	// basic connectivity
	// TODO: Semua yang ada comment, itu yang sebelumnya berjalan dengan baik.
	//       Sampai negara api menyerang.
	// d := GetConfig("database").(map[string]interface{})
	// fmt.Printf("|%+v|", reflect.TypeOf(d["hostname"]))
	// d := viper.Get("database").(map[string]interface{})

	dbconf.Host = viper.GetString("database.hostname") // d["hostname"].(string)
	dbconf.Type = viper.GetString("database.type")     // d["type"].(string)
	// convert dari map viper ke int64 dan convert lagi ke uint16
	// karena int64 terlalu besar untuk menyimpan port yang isinya maksimum hanya 65535
	dbconf.Port = uint16(viper.GetInt64("database.port"))  // uint16(d["port"].(int64))
	dbconf.User = viper.GetString("database.username")     // d["username"].(string)
	dbconf.Pass = viper.GetString("database.password")     // d["password"].(string)
	dbconf.Database = viper.GetString("database.database") // d["database"].(string)

	pak.TryCatchBlock{
		Try: func() {
			Spin.Suffix = " Testing database config"
			succeed, errPing := ora.PingDb(dbconf)
			if !succeed {
				log.Fatalln(errPing)
				os.Exit(3)
			}
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
			Spin.Suffix = " Migrating database"
			if !ora.SetupDb(dbconf) {
				// fmt.Println("Database migration success.")
				fmt.Println("Database migration failed.")
				os.Exit(3)
			}

		},
		Catch: func(e pak.Exception) {
			log.Fatalf("Error raised while running SetupDb: %+v", e)
			os.Exit(3)
		},
		Finally: func() {
			Spin.Suffix = " Database migration success"
		},
	}.Do()

	rt.Dbconf = dbconf

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
	// var modules []*Module

	RunBootLevel1()
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
			Modloc = cwd + "/modules"
		}
	}

	rt.Modloc = Modloc

	Spin.Suffix = " Initiate core modules"
	coreModules, err := ioutil.ReadDir(strings.Join(
		[]string{Modloc, "core"}, "/"))
	pak.ErrHandler(err)

	// Setup router for the first time.
	Routers = SetupRouter()

	for _, coreModule := range coreModules {
		m := LoadCoreModule(coreModule.Name())
		m.Bootstrap()
		m.Router(Routers)
	}
	// Load Core module done

	// fmt.Printf("%+v", modloc)

	// time.Sleep(10 * time.Second)
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
		m := LoadContribModule(contribModule.Name())
		m.Bootstrap()
		m.Router(Routers)
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
	WriteRuntime(rt)
}

func BootstrapAll() {
	Bootstrap(BOOTSTRAP_LEVEL_3)
}
