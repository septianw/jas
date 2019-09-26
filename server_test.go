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
	// "net/http"
	// "net/http/httptest"

	"os"
	"path/filepath"
	"testing"

	pak "github.com/septianw/jas/common"
	// ty "github.com/septianw/jas/types"
	"github.com/stretchr/testify/assert"
)

// func TestBootstrap(t *testing.T) {
// 	BootstrapAll()
// }

// func TestPingRoute(t *testing.T) {
// 	router := SetupRouter()

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/ping", nil)
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, 200, w.Code)
// 	assert.Equal(t, "pong", w.Body.String())
// }

// func TestPingDb(t *testing.T) {
// 	d := ty.Dbconf{
// 		"mysql",
// 		"localhost",
// 		3306,
// 		"asep",
// 		"dummypass",
// 		"ipoint",
// 	}
// 	db := pak.LoadDatabase("/home/asep/gocode/src/github.com/septianw/jas/libs/database.so", d)
// 	connected, err := db.PingDb(d)
// 	t.Log(connected, err)
// 	if !connected {
// 		t.Fail()
// 	}
// }

// func TestMigrate(t *testing.T) {
// 	d := ty.Dbconf{
// 		"mysql",
// 		"localhost",
// 		3306,
// 		"asep",
// 		"dummypass",
// 		"ipoint",
// 	}
// 	db := pak.LoadDatabase("/home/asep/gocode/src/github.com/septianw/jas/libs/database.so", d)

// 	done := db.Migrate("/home/asep/gocode/src/github.com/septianw/jas/migrations", d)
// 	assert.Equal(t, true, done)
// }

// func TestBootstrapLevel0(t *testing.T) {
// 	Bootstrap(BOOTSTRAP_LEVEL_0)

// 	rt := pak.ReadRuntime()
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		t.Fail()
// 	}

// 	t.Logf(`
// 	ListenAddr: %s,
// 	Libloc: %s,
// 	runtime: %+v
// 	`, ListenAddr, Libloc, rt)

// 	t.Logf("\npath: %s", filepath.Join(cwd, "now"))

// 	assert.Equal(t, "192.168.122.1:4519", ListenAddr)
// 	assert.Equal(t, filepath.Join(cwd, "libs"), Libloc)
// 	// assert.IsType(t, ty.Runtime, rt)
// 	// assert.ObjectsAreEqual(ty.Runtime, rt)
// }

// func TestBootstrapLevel1(t *testing.T) {
// 	Bootstrap(BOOTSTRAP_LEVEL_1)

// 	rt := pak.ReadRuntime()
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		t.Fail()
// 	}

// 	t.Logf(`
// 	ListenAddr: %s,
// 	Libloc: %s,
// 	Migloc: %s
// 	runtime: %+v
// 	`, ListenAddr, Libloc, Migloc, rt)

// 	assert.Equal(t, filepath.Join(cwd, "migrations"), Migloc)
// }

func TestBootstrapLevel2(t *testing.T) {
	os.Setenv("cwd", "/home/asep/gocode/src/github.com/septianw/jas")

	Bootstrap(BOOTSTRAP_LEVEL_3)

	rt := pak.ReadRuntime()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fail()
	}

	t.Logf(`
	ListenAddr: %s,
	Libloc: %s,
	Migloc: %s,
	Modloc: %s,
	runtime: %+v
	`, ListenAddr, Libloc, Migloc, Modloc, rt)

	assert.Equal(t, filepath.Join(cwd, "modules"), Modloc)
}

func TestIsModuleEnabled(t *testing.T) {
	os.Setenv("cwd", "/home/asep/gocode/src/github.com/septianw/jas")
	isUser := IsModuleEnabled("user")
	if isUser == false {
		t.Fail()
	}

	isNotFound := IsModuleEnabled("notfound")
	if isNotFound == true {
		t.Fail()
	}

	t.Logf("is User: %+v", isUser)
	t.Logf("is NotFound: %+v", isNotFound)
}
