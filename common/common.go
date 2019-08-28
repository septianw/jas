package common

import (
	"encoding/gob"
	"log"
	"os"
	"strings"

	// "github.com/juju/loggo"
	"runtime/debug"

	ty "github.com/septianw/jas/types"
)

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

func ReadRuntime() ty.Runtime {
	var out ty.Runtime

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
