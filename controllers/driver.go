package controllers

import (
	"fmt"
	"strconv"

	"github.com/electra-systems/athena/storage"
	"github.com/electra-systems/athena/utils"

	"github.com/go-redis/redis"
	"github.com/uber/h3-go"
)

type DriverController struct {
	DB storage.StorageInstance
}

type Response struct {
	Message string
	Err     error
	Data    map[string]interface{}
}

type DriverLocationData struct {
	DriverId string `json:"driver_id"`
	Lat      string `json:"lat"`
	Lng      string `json:"lng"`
}

func (c *DriverController) IndexLocation(data DriverLocationData) Response {

	lat, _ := strconv.ParseFloat(data.Lat, 64)
	lng, _ := strconv.ParseFloat(data.Lng, 64)

	h3Index := utils.IndexLatLng(h3.GeoCoord{Latitude: lat, Longitude: lng})
	stringifiedIndex := utils.H3IndexToString(h3Index)

	lastDriverLocationIndex, err := c.DB.Driver.Get(data.DriverId)

	if err != redis.Nil && err != nil {

		return Response{
			Message: "Last driver location lookup Err",
			Err:     err,
		}

	}

	fmt.Println(stringifiedIndex, lastDriverLocationIndex)

	if stringifiedIndex == lastDriverLocationIndex {
		return Response{
			Message: "Driver hasn't changed position",
			Err:     nil,
		}
	}

	_, err = c.DB.Driver.Set(data.DriverId, uint64(h3Index))

	if err != nil {

		return Response{
			Message: "Updating driver location failed",
			Err:     err,
		}

	}

	_, err = c.DB.Car.RemoveFromList(lastDriverLocationIndex, data.DriverId)

	if err != nil {

		return Response{
			Message: "Updating old index failed",
			Err:     err,
		}

	}

	_, err = c.DB.Car.InsertIntoList(stringifiedIndex, data.DriverId)

	if err != nil {
		return Response{
			Message: "Updating new index failed",
			Err:     err,
		}
	}

	reponseValue := map[string]interface{}{
		"driver_id":           data.DriverId,
		"last_driver_index":   lastDriverLocationIndex,
		"latest_driver_index": stringifiedIndex,
		"lat":                 lat,
		"lng":                 lng,
	}

	return Response{
		Data: reponseValue,
	}

}
