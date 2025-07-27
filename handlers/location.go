package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sch8ill/masstrack/db"
	"github.com/sch8ill/masstrack/location"
)

func Locations(db *db.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var locations []location.Location
		var err error

		device := c.Query("device")
		if device != "" {
			locations, err = db.DeviceLocations(device)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		} else {
			locations, err = db.CurrentLocations()
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}

		c.JSON(http.StatusOK, locations)
	}
}
