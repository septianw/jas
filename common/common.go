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
	"fmt"
	"net/http"
	"plugin"
	"runtime/debug"

	// pak "github.com/septianw/jas/common"

	"github.com/gin-gonic/gin"
	ty "github.com/septianw/jas/types"
)

type Exception interface{}

type TryCatchBlock struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

/*
ERROR CODE LEGEND:
error containt 4 digits,
first digit represent error location either module or main app
1 for main app
2 for module

second digit represent error at level app or database
1 for app
2 for database

third digit represent error with input variable or variable manipulation
0 for skipping this error
1 for input validation error
2 for variable manipulation error
3 for input query result empty record

fourth digit represent error with logic, this type of error have
increasing error number based on which part of code that error.
0 for skipping this error
1 for unknown logical error
2 for whole operation fail, operation end unexpectedly
3 for whole operation fail, access forbidden
*/

const DATABASE_EXEC_FAIL_CODE = 2200
const MODULE_OPERATION_FAIL_CODE = 2102
const INPUT_VALIDATION_FAIL_CODE = 2110
const RECORD_NOT_FOUND_CODE = 2230
const PAGE_NOT_FOUND_CODE = 2100
const NOT_ACCEPTABLE_CODE = 2112
const FORBIDDEN_CODE = 2103
const UNKNOWN_ERROR_CODE = 2101

var NOT_ACCEPTABLE = gin.H{"code": NOT_ACCEPTABLE_CODE, "message": "You are trying to request something not acceptible here."}
var PAGE_NOT_FOUND = gin.H{"code": PAGE_NOT_FOUND_CODE, "message": "You are find something we can't found it here."}
var RECORD_NOT_FOUND = gin.H{"code": RECORD_NOT_FOUND_CODE, "message": "You are find something we can't found it here."}

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

func SendHttpError(c *gin.Context, errType uint, err error) {
	// FIXME: kurang yang forbidden
	switch errType {
	case DATABASE_EXEC_FAIL_CODE:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": DATABASE_EXEC_FAIL_CODE,
			"message": fmt.Sprintf("DATABASE_EXEC_FAIL: %s", err.Error())})
		break
	case INPUT_VALIDATION_FAIL_CODE:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": INPUT_VALIDATION_FAIL_CODE,
			"message": fmt.Sprintf("INPUT_VALIDATION_FAIL: %s", err.Error())})
		break
	case PAGE_NOT_FOUND_CODE:
		c.AbortWithStatusJSON(http.StatusNotFound, PAGE_NOT_FOUND)
		break
	case RECORD_NOT_FOUND_CODE:
		c.AbortWithStatusJSON(http.StatusNotFound, RECORD_NOT_FOUND)
		break
	case FORBIDDEN_CODE:
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": FORBIDDEN_CODE,
			"message": fmt.Sprintf("FORBIDDEN: %s", err.Error())})
		break
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"code": errType,
				"message": fmt.Sprintf("UNKNOWN_ERROR: %s", err.Error())})
	}
}
