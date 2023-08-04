package controllers

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uncleDecart/gha-stat-collector/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const PerPageDefault = 25

type dateQueryParam struct {
	Eq   string
	Date time.Time
}

func parseDateQueryParam(q string) (dateQueryParam, error) {
	time.Now()
	r := regexp.MustCompile("(lt|lte|eq|gte|gt)((\\d{7,}))")
	if !r.MatchString(q) {
		return dateQueryParam{}, fmt.Errorf("Query parameter %s has invalid format", q)
	}

	match := r.FindStringSubmatch(q)

	timestamp, err := strconv.ParseInt(match[2], 10, 64)
	if err != nil {
		return dateQueryParam{}, err
	}

	return dateQueryParam{Eq: match[1], Date: time.Unix(timestamp, 0)}, nil
}

func ComposeFilterFromQuery(c *gin.Context) bson.M {
	filter := bson.M{}

	if name, ok := c.GetQuery("name"); ok {
		filter["name"] = name
	}

	if start, ok := c.GetQuery("start"); ok {
		dqp, err := parseDateQueryParam(start)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		// Since API parameters are compatible with MongoDb which we are using
		filter["start"] = bson.M{"$" + dqp.Eq: strconv.FormatInt(dqp.Date.Unix(), 10)}
	}

	if end, ok := c.GetQuery("end"); ok {
		dqp, err := parseDateQueryParam(end)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		// Since API parameters are compatible with MongoDb which we are using
		filter["end"] = bson.M{"$" + dqp.Eq: strconv.FormatInt(dqp.Date.Unix(), 10)}
	}

	if outcome, ok := c.GetQuery("outcome"); ok {
		filter["outcome"] = strings.ToLower(outcome)
	}

	if arch, ok := c.GetQuery("arch"); ok {
		filter["arch"] = arch
	}

	return filter
}

func Read(c *gin.Context) {
	filter := ComposeFilterFromQuery(c)

	perPage := int64(PerPageDefault)
	if pp, ok := c.GetQuery("perpage"); ok {
		intPp, err := strconv.ParseInt(pp, 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		perPage = intPp
	}
	curPage := int64(1)
	if cp, ok := c.GetQuery("curpage"); ok {
		intCp, err := strconv.ParseInt(cp, 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		curPage = intCp
	}
	skip := (curPage - 1) * perPage

	options := options.Find().SetSkip(int64(skip)).SetLimit(int64(perPage))

	client, err := models.GetClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Disconnect(context.TODO())
	collection := client.Database(DbName).Collection(CollectionName)

	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	var data []ActionLogEntry

	for cursor.Next(context.TODO()) {
		var el ActionLogEntry
		if err := cursor.Decode(&el); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		data = append(data, el)
	}

	totalDocuments, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := ActionLogEntrySearch{
		TotalPages: int64(math.Ceil(float64(totalDocuments) / float64(perPage))),
		PerPage:    perPage,
		CurPage:    curPage,
		Data:       data,
	}

	c.JSON(http.StatusOK, response)
}
