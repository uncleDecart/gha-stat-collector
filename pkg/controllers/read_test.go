package controllers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/uncleDecart/gha-stat-collector/pkg/controllers"
	"go.mongodb.org/mongo-driver/bson"
	"gotest.tools/assert"
)

func TestComposeFilterFromQueryName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Name
	url := "?name=testoname&start=eq1234567&end=lte1234567&successful=true&arch=arm"
	c.Request, _ = http.NewRequest("GET", url, nil)
	expected := bson.M{
		"name":       "testoname",
		"start":      bson.M{"$eq": int64(1234567)},
		"end":        bson.M{"$lte": int64(1234567)},
		"successful": bson.M{"$ne": false},
		"arch":       "arm",
	}

	got := controllers.ComposeFilterFromQuery(c)
	fmt.Println(got)
	fmt.Println("----")
	fmt.Println(expected)
	assert.Equal(t, true, reflect.DeepEqual(got, expected))
}
