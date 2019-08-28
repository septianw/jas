package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ty "github.com/septianw/jas/types"
	"github.com/stretchr/testify/assert"
)

// func TestBootstrap(t *testing.T) {
// 	BootstrapAll()
// }

func TestPingRoute(t *testing.T) {
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestPingDb(t *testing.T) {
	d := ty.Dbconf{
		"mysql",
		"localhost",
		3306,
		"asep",
		"dummypass",
		"ipoint",
	}
	db := LoadDatabase("/home/asep/gocode/src/github.com/septianw/jas/libs/database.so", d)
	connected, err := db.PingDb(d)
	t.Log(connected, err)
	if !connected {
		t.Fail()
	}
}

func TestMigrate(t *testing.T) {
	d := ty.Dbconf{
		"mysql",
		"localhost",
		3306,
		"asep",
		"dummypass",
		"ipoint",
	}
	db := LoadDatabase("/home/asep/gocode/src/github.com/septianw/jas/libs/database.so", d)

	db.SetupDb(d)
}
