package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uncleDecart/gha-stat-collector/pkg/models"
)

type StepLogEntry struct {
	Id          uint64 `json:"id" binding:"required"`
	ExecTime    string `json:"exec_time" binding:"required"`
	Successfull bool   `json:"successful" binding:"required"`
}

type ActionLogEntry struct {
	Name        string         `json:"name" binding:"required"`
	Start       string         `json:"start" binding:"required"`
	End         string         `json:"end" binding:"required"`
	ExecTime    string         `json:"exec_time" binding:"required"`
	Successfull bool           `json:"successful" binding:"required"`
	Arch        string         `json:"arch" binding:"required"`
	Steps       []StepLogEntry `json:"steps"`
}

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

	collection := client.Database("records").Collection("action-logs")
	insertResult, err := collection.InsertOne(context.TODO(), entry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": fmt.Sprint("Inserted document ", insertResult.InsertedID),
	})
}
