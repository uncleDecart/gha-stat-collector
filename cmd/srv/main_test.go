package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/uncleDecart/gha-stat-collector/pkg/controllers"
	"gotest.tools/assert"
)

func TestEnvVariables(t *testing.T) {
	// All variables from .env file must be present
	_, present := os.LookupEnv("ACCESS_TOKEN")
	assert.Equal(t, true, present)

	_, present = os.LookupEnv("MONGO_USERNAME")
	assert.Equal(t, true, present)

	_, present = os.LookupEnv("MONGO_PASSWORD")
	assert.Equal(t, true, present)

	_, present = os.LookupEnv("MONGO_URL")
	assert.Equal(t, true, present)
}

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
	at := os.Getenv("ACCESS_TOKEN")

	router := setupRouter()

	val := true

	body := controllers.ActionLogEntry{
		Name:        "test-action",
		Start:       strconv.FormatInt(time.Now().Unix(), 10),
		End:         strconv.FormatInt(time.Now().Unix(), 10),
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

func TestGetTimingRoute(t *testing.T) {
	at := os.Getenv("ACCESS_TOKEN")

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/timing", nil)
	req.Header.Set("auth", at)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
