package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/uber/h3-go"
)

type Coord struct {
	Lat, Lng float64
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func IndexLatLng(coordinates h3.GeoCoord) h3.H3Index {
	return h3.FromGeo(coordinates, 8)
}

func H3IndexToString(index h3.H3Index) string {
	return fmt.Sprintf("%v", index)
}

func FormatH3Index(index h3.H3Index) string {
	return fmt.Sprintf("%#x\n", index)
}
