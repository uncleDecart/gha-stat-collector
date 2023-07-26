package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/uncleDecart/gha-stat-collector/pkg/controllers"
	"gotest.tools/assert"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	type PingResponse struct {
		Message string `json:"message"`
	}
	var got PingResponse
	json.NewDecoder(w.Body).Decode(&got)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, PingResponse{Message: "pong"}, got)
}

func TestPostTimingRoute(t *testing.T) {
	at := "qwerty"
	os.Setenv("ACCESS_TOKEN", at)

	router := setupRouter()

	val := true

	body := controllers.ActionLogEntry{
		Name:        "test-action",
		Start:       "2023-02-01",
		End:         "2023-02-05",
		Successfull: &val,
		Arch:        "Arm",
	}

	b, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/timing", bytes.NewReader(b))
	req.Header.Set("auth", at)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
