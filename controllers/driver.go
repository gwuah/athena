package controllers

import (
	"fmt"
	"log"
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

func (c *DriverController) FindClosestDrivers(data DriverLocationData, neighbours int) Response {
	lat, _ := strconv.ParseFloat(data.Lat, 64)
	lng, _ := strconv.ParseFloat(data.Lng, 64)

	h3Index := utils.IndexLatLng(h3.GeoCoord{Latitude: lat, Longitude: lng})

	rings := h3.KRing(h3Index, neighbours)

	cars := []string{}

	for _, value := range rings {
		matchedCars, err := c.DB.Car.All(utils.H3IndexToString(value))

		if err != nil {
			log.Println("Failed To Retrieve Active Drivers")
			continue
		}

		cars = append(cars, matchedCars...)
	}

	reponseValue := map[string]interface{}{
		"polygons": []interface{}{utils.GeneratePolygons(rings)},
		"drivers":  cars,
	}

	return Response{
		Data: reponseValue,
	}
}
