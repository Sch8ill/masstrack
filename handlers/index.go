package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"currentTime": time.Now().Format(timeFormat),
	})
}
