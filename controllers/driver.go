package controllers

import (
	"github.com/electra-systems/athena/services"
	"github.com/electra-systems/athena/storage"
	"github.com/electra-systems/athena/utils"
	"github.com/gin-gonic/gin"
	"github.com/uber/h3-go"
)

type DriverController struct {
	DB storage.StorageInstance
}

func (d *DriverController) IndexLocation(c *gin.Context) {

	var data services.Payload

	if c.BindJSON(&data) != nil {
		c.JSON(500, gin.H{
			"message": "Error",
		})
		return
	}

	dbMap := map[string]storage.Redis{
		"driver":      d.DB.Driver,
		"supplyIndex": d.DB.Car,
	}

	geoIndexer := services.NewGeoIndex(dbMap, "driver", "supplyIndex")

	response, err := geoIndexer.Index(data)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed To Index",
			"error":   err,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Successfully Indexed Coordinates",
		"data":    response,
	})

}

func (d *DriverController) FindClosestDrivers(c *gin.Context) {

	var data services.Payload

	if c.BindJSON(&data) != nil {
		c.JSON(500, gin.H{
			"message": "Error",
		})
		return
	}

	dbMap := map[string]storage.Redis{
		"driver":      d.DB.Driver,
		"supplyIndex": d.DB.Car,
	}

	searchService := services.NewSearch(dbMap, "driver", "supplyIndex")

	response, err := searchService.Closest(data, 2)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to search",
			"error":   err,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Closest Drivers Retrieved",
		"data":    response,
	})

}

func (d *DriverController) GetMapOverlay(c *gin.Context) {
	var data services.Payload

	if c.BindJSON(&data) != nil {
		c.JSON(500, gin.H{
			"message": "Error",
		})
		return
	}

	var parsedValue = utils.IndexCoordinates(utils.IndexCoordinatesProps{
		Lat: data.Lat,
		Lng: data.Lng,
	})

	rings := h3.KRing(parsedValue.Index, 2)

	c.JSON(200, gin.H{
		"message": "Overlay Displayed",
		"view":    utils.GeneratePolygons(rings),
	})

}
