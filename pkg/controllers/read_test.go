package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/uncleDecart/gha-stat-collector/pkg/controllers"
	"go.mongodb.org/mongo-driver/bson"
	"gotest.tools/assert"
)

// CompareBsonM compares two bson.M objects.
func CompareBsonM(a, b bson.M) bool {
	if len(a) != len(b) {
		return false
	}

	keysA := make([]string, 0, len(a))
	keysB := make([]string, 0, len(b))

	// Get all keys from both maps
	for key := range a {
		keysA = append(keysA, key)
	}
	for key := range b {
		keysB = append(keysB, key)
	}

	// Sort the keys to ensure consistent ordering for comparison
	sort.Strings(keysA)
	sort.Strings(keysB)

	// Check if keys match
	if !reflect.DeepEqual(keysA, keysB) {
		return false
	}

	// Check if values match for each key
	for _, key := range keysA {
		if !reflect.DeepEqual(a[key], b[key]) {
			return false
		}
	}

	return true
}

func TestComposeFilterFromQueryName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	url := "?name=testoname&outcome=success&arch=arm"
	c.Request, _ = http.NewRequest("GET", url, nil)
	expected := bson.M{
		"name":    "testoname",
		"outcome": "success",
		"arch":    "arm",
	}

	got := controllers.ComposeFilterFromQuery(c)
	assert.Equal(t, true, reflect.DeepEqual(got, expected))
}
