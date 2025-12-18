package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sch8ill/masstrack/db"
	"github.com/sch8ill/masstrack/location"
)

const timeFormat string = "2006-01-02T15:04"

func Locations(db *db.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var locations []location.Location
		var err error

		timespan := true
		start, err := time.Parse(timeFormat, c.Query("start"))
		if err != nil {
			timespan = false
		}

		end, err := time.Parse(timeFormat, c.Query("end"))
		if err != nil {
			timespan = false
		}

		device := c.Query("device")
		if device != "" {
			locations, err = db.DeviceLocations(device)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusOK, locations)
			return
		}

		if timespan {
			locations, err = db.TimespanLocations(start, end)
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
