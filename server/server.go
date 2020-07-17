package server

import (
	"fmt"
	"log"

	"github.com/electra-systems/athena/controllers"
	"github.com/electra-systems/athena/storage"
	"github.com/electra-systems/athena/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleElectronWebsocketConnection(h *Hub) func(c *gin.Context) {
	return func(c *gin.Context) {
		electronId := c.Param("id")
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		fmt.Println("Connection Recieved from", electronId)

		if err != nil {
			log.Println("Failed to setup websocket conn ..")
			log.Println(err)
			return
		}

		electron := &Electron{hub: h, send: make(chan []byte), conn: conn, id: electronId}

		go electron.readMessages()
		go electron.writeMessagesToClient()
	}
}

func Init(db storage.StorageInstance) {

	driverController := controllers.DriverController{DB: db}

	r := gin.Default()

	hub := newHub()

	go hub.Init()

	r.Use(utils.CORSMiddleware())

	r.GET("/electron-ws/:id", handleElectronWebsocketConnection(hub))

	r.POST("/index-driver-location", driverController.IndexLocation)

	r.POST("/closest-drivers", driverController.FindClosestDrivers)

	r.POST("/get-map-overlay", driverController.GetMapOverlay)

	r.Run()

}
