package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uncleDecart/gha-stat-collector/pkg/models"
)

func Publish(c *gin.Context) {
	var entry ActionLogEntry

	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := models.GetClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database(DB_NAME).Collection(COLLECTION_NAME)
	insertResult, err := collection.InsertOne(context.TODO(), entry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": fmt.Sprint("Inserted document ", insertResult.InsertedID),
	})
}
