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

	r.POST("/index-driver-location", driverController.IndexLocation)

	r.POST("/closest-drivers", driverController.FindClosestDrivers)

	r.POST("/get-map-overlay", driverController.GetMapOverlay)

	r.Run()

}
