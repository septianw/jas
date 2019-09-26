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
package common

import (
	"encoding/gob"
	"log"
	"os"
	"strings"

	// "github.com/juju/loggo"
	"plugin"
	"runtime/debug"

	// pak "github.com/septianw/jas/common"

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

// This function will load *.so library without parsing its function.
// After load library with this function you need to lookup your function.
func LoadSo(path string) *plugin.Plugin {
	plug, err := plugin.Open(path)
	ErrHandler(err)

	return plug
}

func ReadRuntime() ty.Runtime {
	var out ty.Runtime

	RuntimeFile, err := os.OpenFile("/tmp/shinyRuntimeFile", os.O_RDWR, 0600)
	ErrHandler(err)

	dec := gob.NewDecoder(RuntimeFile)
	err = dec.Decode(&out)
	ErrHandler(err)

	err = RuntimeFile.Close()
	ErrHandler(err)

	return out
}

func WriteRuntime(rt ty.Runtime) {
	RuntimeFile, err := os.OpenFile("/tmp/shinyRuntimeFile", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	ErrHandler(err)

	enc := gob.NewEncoder(RuntimeFile)
	err = enc.Encode(rt)
	ErrHandler(err)

	err = RuntimeFile.Close()
	ErrHandler(err)
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

func LoadDatabase(libpath string, d ty.Dbconf) ty.Database {
	// rt := pak.ReadRuntime()

	// pak.ErrHandler(errors.New(rt.Libloc))
	// pak.ErrHandler(errors.New(filepath.Join(rt.Libloc, "database.so")))

	plug := LoadSo(libpath)
	symd, err := plug.Lookup("Database")
	ErrHandler(err)
	sd := symd.(ty.Database)

	return sd
}
