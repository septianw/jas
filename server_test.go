package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestPingDb(t *testing.T) {
	d := Dbconf{
		"mysql",
		"localhost",
		3306,
		"asep",
		"dummypass",
		"ipoint",
	}
	connected, err := PingDb(d)
	t.Log(connected, err)
	if !connected {
		t.Fail()
	}
}

func TestMigrate(t *testing.T) {
	d := Dbconf{
		"mysql",
		"localhost",
		3306,
		"asep",
		"dummypass",
		"ipoint",
	}

	SetupDb(d)
}
