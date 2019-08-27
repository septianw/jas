package common

import (
	"database/sql"
	"encoding/gob"
	"log"
	"os"
	"strings"

	// "github.com/juju/loggo"
	"runtime/debug"
)

type Dbconf struct {
	Type     string
	Host     string
	Port     uint16
	User     string
	Pass     string
	Database string
}

type OpenDbFunc func(d Dbconf) (*sql.DB, error)

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

type Status struct {
	Installed uint8
}

type Exception interface{}

type TryCatchBlock struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

func (tc TryCatchBlock) Do() {
	if tc.Finally != nil {
		defer tc.Finally()
	}
	if tc.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				ErrHandler(r.(error))
				tc.Catch(r)
			}
		}()
	}
	tc.Try()
}

func ReadRuntime() Runtime {
	var out Runtime

	RuntimeFile, err := os.OpenFile("/tmp/shinyRuntimeFile", os.O_RDWR|os.O_CREATE, 0400)
	ErrHandler(err)

	dec := gob.NewDecoder(RuntimeFile)
	err = dec.Decode(&out)
	ErrHandler(err)

	err = RuntimeFile.Close()
	ErrHandler(err)

	return out
}

func ErrHandler(err error) {
	var ginMode = os.Getenv("GIN_MODE")
	var stage = os.Getenv("STAGE")
	// var logger = loggo.GetLogger("")

	if err != nil {
		if (strings.Compare(ginMode, "release") == 0) ||
			(strings.Compare(stage, "release") == 0) {
			log.Printf("Eew, error occured: %+v", err)
			// } else if (strings.Compare(ginMode, "development") == 0) ||
			// 	(strings.Compare(stage, "development") == 0) {
			// 	log.Printf("Thing to be done: %+v", err)
			// 	debug.PrintStack()
			// logger.Debugf("\nError occured: %+v\n", err)
		} else {
			log.Printf("Thing to be done: %+v", err)
			debug.PrintStack()
			log.Println(err)
		}
	}
}
