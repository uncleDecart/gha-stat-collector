package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	assert.Equal(t, true, present, "ACCESS_TOKEN variable should be present")

	_, present = os.LookupEnv("MONGO_USERNAME")
	assert.Equal(t, true, present, "MONGO_USERNAME variable should be present")

	_, present = os.LookupEnv("MONGO_PASSWORD")
	assert.Equal(t, true, present, "MONGO_PASSWORD variable should be present")

	_, present = os.LookupEnv("MONGO_URL")
	assert.Equal(t, true, present, "MONGO_URL variable should be present")
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

	body := controllers.ActionLogEntry{
		Name:    "test-action",
		Start:   strconv.FormatInt(time.Now().Unix(), 10),
		End:     strconv.FormatInt(time.Now().Unix(), 10),
		Outcome: "success",
		Arch:    "Arm",
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

func TestGetTimingRouteFilter(t *testing.T) {
	// Warning: keep in mind that previous entries are not deleted
	at := os.Getenv("ACCESS_TOKEN")
	router := setupRouter()
	w := httptest.NewRecorder()

	action1 := controllers.ActionLogEntry{
		Name:    "testfilternameaction1",
		Start:   strconv.FormatInt(time.Now().Unix(), 10),
		End:     strconv.FormatInt(time.Now().Unix()+10, 10),
		Outcome: "success",
		Arch:    "Arm",
	}
	action2 := controllers.ActionLogEntry{
		Name:    "testfilternameaction2",
		Start:   strconv.FormatInt(time.Now().Unix()+20, 10),
		End:     strconv.FormatInt(time.Now().Unix()+30, 10),
		Outcome: "failure",
		Arch:    "x86_64",
	}

	actions := []controllers.ActionLogEntry{action1, action2}

	for _, actionToSend := range actions {
		b, _ := json.Marshal(actionToSend)

		req, _ := http.NewRequest("POST", "/api/v1/timing", bytes.NewReader(b))
		req.Header.Set("auth", at)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		// Skip responses from insert requests
		io.ReadAll(w.Body) // response body is []byte
	}

	// Test name filtering
	url := fmt.Sprintf("/api/v1/timing?name=%s", action1.Name)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("auth", at)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	body, err := io.ReadAll(w.Body) // response body is []byte
	assert.Equal(t, nil, err, "Failed to read response body")
	var result controllers.ActionLogEntrySearch
	err = json.Unmarshal(body, &result)
	assert.Equal(t, nil, err, "Failed to unmarshal response body to struct")
	expected := []controllers.ActionLogEntry{action1}
	assert.Equal(t, true, controllers.CompareActionLogEntrySlice(expected, result.Data))

	// Test start filtering
	url = fmt.Sprintf("/api/v1/timing?start=eq%s", action2.Start)
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("auth", at)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	body, err = io.ReadAll(w.Body) // response body is []byte
	assert.Equal(t, nil, err, "Failed to read response body")
	err = json.Unmarshal(body, &result)
	assert.Equal(t, nil, err, "Failed to unmarshal response body to struct")
	expected = []controllers.ActionLogEntry{action2}
	assert.Equal(t, true, controllers.CompareActionLogEntrySlice(expected, result.Data))

	// Test end filtering
	url = fmt.Sprintf("/api/v1/timing?end=leq%s", action1.End)
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("auth", at)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	body, err = io.ReadAll(w.Body)
	assert.Equal(t, nil, err, "Failed to read response body")
	err = json.Unmarshal(body, &result)
	assert.Equal(t, nil, err, "Failed to unmarshal response body to struct")
	expected = []controllers.ActionLogEntry{action1}
	assert.Equal(t, true, controllers.CompareActionLogEntrySlice(expected, result.Data))

	// Test outcome filtering; checking false
	url = fmt.Sprintf("/api/v1/timing?outcome=%s", action2.Outcome)
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("auth", at)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	body, err = io.ReadAll(w.Body)
	bodyString := string(body)
	fmt.Println(bodyString)
	assert.Equal(t, nil, err, "Failed to read response body")
	err = json.Unmarshal(body, &result)
	assert.Equal(t, nil, err, "Failed to unmarshal response body to struct")
	expected = []controllers.ActionLogEntry{action2}
	assert.Equal(t, true, controllers.CompareActionLogEntrySlice(expected, result.Data))

	// Test arch filtering
	url = fmt.Sprintf("/api/v1/timing?arch=%s", action2.Arch)
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("auth", at)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	body, err = io.ReadAll(w.Body) // response body is []byte
	assert.Equal(t, nil, err, "Failed to read response body")
	err = json.Unmarshal(body, &result)
	assert.Equal(t, nil, err, "Failed to unmarshal response body to struct")
	expected = []controllers.ActionLogEntry{action2}
	assert.Equal(t, true, controllers.CompareActionLogEntrySlice(expected, result.Data))
}
