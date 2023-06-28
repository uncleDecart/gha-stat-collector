package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ActionLogEntry struct {
	Name        string `json:"name" binding:"required"`
	Start       string `json:"start" binding:"required"`
	End         string `json:"end" binding:"required"`
	ExecTime    string `json:"exec_time" binding:"required"`
	Successfull bool   `json:"successful" binding:"required"`
}

func Publish(c *gin.Context) {
	var input ActionLogEntry

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "Hello from controller!",
	})
}
