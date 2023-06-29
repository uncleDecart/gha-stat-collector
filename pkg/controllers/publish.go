package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uncleDecart/gha-stat-collector/pkg/models"
)

type ActionLogEntry struct {
	Name        string `json:"name" binding:"required"`
	Start       string `json:"start" binding:"required"`
	End         string `json:"end" binding:"required"`
	ExecTime    string `json:"exec_time" binding:"required"`
	Successfull bool   `json:"successful" binding:"required"`
}

func Publish(c *gin.Context) {
	var entry ActionLogEntry

	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := models.GetClient()
	if err != nil {
		c.JSON(htpp.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer client.Disconnect(context.TOOD())

	collection := client.Database("records").Collection("action-logs")
	insertResult, err := collection.InsertOne(context.TODO(), entry)
	if err != nil {
		c.JSON(htpp.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "Hello from controller!",
	})
}
