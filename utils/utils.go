package utils

import (
	"fmt"
	"strconv"

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

func H3ToPolyline(h3idx h3.H3Index) []Coord {
	hexBoundary := h3.ToGeoBoundary(h3idx)
	hexBoundary = append(hexBoundary, hexBoundary[0])

	arr := []Coord{}

	for _, value := range hexBoundary {
		arr = append(arr, Coord{Lat: value.Latitude, Lng: value.Longitude})
	}

	return arr
}

func GeneratePolygons(rings []h3.H3Index) [][]Coord {
	arr := [][]Coord{}

	for _, value := range rings {
		arr = append(arr, H3ToPolyline(value))
	}

	return arr
}

// type ParseAndIndexProps struct {
// 	Lat        string
// 	Lng        string
// 	Neighbours int
// }

// type ParseAndIndexReturnValue struct {
// 	Rings []h3.H3Index
// }

type IndexCoordinatesProps struct {
	Lat string
	Lng string
}

type IndexCoordinatesReturnValue struct {
	Lat, Lng float64
	Index    h3.H3Index
}

func IndexCoordinates(props IndexCoordinatesProps) IndexCoordinatesReturnValue {
	lat, _ := strconv.ParseFloat(props.Lat, 64)
	lng, _ := strconv.ParseFloat(props.Lng, 64)

	return IndexCoordinatesReturnValue{
		Index: IndexLatLng(h3.GeoCoord{Latitude: lat, Longitude: lng}),
		Lat:   lat,
		Lng:   lng,
	}
}

func StringifyLngLat(props h3.GeoCoord) string {

	return "" + fmt.Sprintf("%f", props.Longitude) + "," + fmt.Sprintf("%f", props.Latitude)

}
