package services

import (
	"github.com/electra-systems/athena/controllers"
	"github.com/electra-systems/athena/storage"
	Utils "github.com/electra-systems/athena/utils"

	"github.com/gin-gonic/gin"
)

func Init(db storage.StorageInstance) {
	driverController := controllers.DriverController{DB: db}

	r := gin.Default()

	r.Use(Utils.CORSMiddleware())

	r.POST("/index-driver-location", func(c *gin.Context) {
		var data controllers.DriverLocationData

		if c.BindJSON(&data) != nil {

			c.JSON(500, gin.H{
				"message": "Error",
			})

			return
		}

		response := driverController.IndexLocation(data)

		if response.Err != nil {
			c.JSON(500, gin.H{
				"message": response.Message,
				"error":   response.Err,
			})
			return
		}

		c.JSON(200, gin.H{
			"message": response.Message,
			"data":    response.Data,
		})

	})

	r.GET("/get-closest-drivers", func(c *gin.Context) {
		var data controllers.DriverLocationData

		if c.BindJSON(&data) != nil {

			c.JSON(500, gin.H{
				"message": "Error",
			})

			return
		}

		response := driverController.FindClosestDrivers(data, 0)

		if response.Err != nil {
			c.JSON(500, gin.H{
				"message": response.Message,
				"error":   response.Err,
			})
			return
		}

		c.JSON(200, gin.H{
			"message": response.Message,
			"data":    response.Data,
		})

	})

	r.Run()

}
