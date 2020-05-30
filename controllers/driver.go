package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/electra-systems/athena/storage"
	"github.com/electra-systems/athena/utils"
	"github.com/go-redis/redis"
	"github.com/gwuah/scully"

	"github.com/uber/h3-go"
)

type DriverController struct {
	DB storage.StorageInstance
}

type DistanceAndTime struct {
	Time     float64 `json:"time"`
	Distance float64 `json:"distance"`
}

type DriverWithTimeAndDistance struct {
	Driver DriverInstance  `json:"driver"`
	DT     DistanceAndTime `json:"dt"`
}

type Response struct {
	Message string
	Err     error
	Data    map[string]interface{}
}

type DriverLocationData struct {
	Id  string `json:"id"`
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type TripRequestData struct {
	Id        string `json:"id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type GeoCoord struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type DriverInstance struct {
	Id             string   `json:"id"`
	LastKnownIndex string   `json:"lastKnownIndex"`
	Coordinates    GeoCoord `json:"coordinates"`
}

func (c *DriverController) IndexLocation(data DriverLocationData) Response {

	indexedValue := utils.IndexCoordinates(utils.IndexCoordinatesProps{
		Lat: data.Lat,
		Lng: data.Lng,
	})

	stringifiedIndex := utils.H3IndexToString(indexedValue.Index)

	storedDriverData, err := c.DB.Driver.Get(data.Id)

	if err != redis.Nil && err != nil {
		return Response{
			Message: "Last driver location lookup Err",
			Err:     err,
		}
	}

	var driverInstance DriverInstance

	err = json.Unmarshal([]byte(storedDriverData), &driverInstance)

	fmt.Println(stringifiedIndex, driverInstance.LastKnownIndex)

	if stringifiedIndex == driverInstance.LastKnownIndex {
		return Response{
			Message: "Driver hasn't changed position",
			Err:     nil,
		}
	}

	instance := DriverInstance{Id: data.Id, LastKnownIndex: stringifiedIndex, Coordinates: GeoCoord{
		Latitude:  indexedValue.Lat,
		Longitude: indexedValue.Lng,
	}}

	marshalledValue, err := json.Marshal(instance)

	fmt.Println(string(marshalledValue))

	_, err = c.DB.Driver.Set(data.Id, marshalledValue)

	if err != nil {
		return Response{
			Message: "Updating driver location failed",
			Err:     err,
		}
	}

	_, err = c.DB.Car.RemoveFromList(driverInstance.LastKnownIndex, data.Id)

	if err != nil {
		return Response{
			Message: "Updating old index failed",
			Err:     err,
		}
	}

	_, err = c.DB.Car.InsertIntoList(stringifiedIndex, data.Id)

	if err != nil {
		return Response{
			Message: "Updating new index failed",
			Err:     err,
		}
	}

	reponseValue := map[string]interface{}{
		"driver_id":           data.Id,
		"last_driver_index":   driverInstance.LastKnownIndex,
		"latest_driver_index": stringifiedIndex,
		"coordinates": map[string]interface{}{
			"latitude":  indexedValue.Lat,
			"longitude": indexedValue.Lng,
		},
	}

	return Response{
		Data:    reponseValue,
		Message: "Success",
	}

}

func (c *DriverController) GetMapOverlay(data DriverLocationData, neighbours int) Response {
	var parsedValue = utils.IndexCoordinates(utils.IndexCoordinatesProps{
		Lat: data.Lat,
		Lng: data.Lng,
	})

	rings := h3.KRing(parsedValue.Index, neighbours)

	return Response{
		Data: map[string]interface{}{
			"view": utils.GeneratePolygons(rings),
		},
	}
}

func (c *DriverController) FindClosestDrivers(data TripRequestData, neighbours int) Response {

	var parsedValue = utils.IndexCoordinates(utils.IndexCoordinatesProps{
		Lat: data.Latitude,
		Lng: data.Longitude,
	})

	rings := h3.KRing(parsedValue.Index, neighbours)

	cars := []interface{}{}

	for _, value := range rings {
		matchedCars, err := c.DB.Car.All(utils.H3IndexToString(value))

		if err != nil {
			log.Println("Failed To Retrieve Active Drivers")
			continue
		}

		if len(matchedCars) == 0 {
			// if we pass an empty array to Mget, it throws an error
			// so we just exit the current iteration and move on
			continue
		}

		driverDetails, err := c.DB.Driver.MGet(matchedCars)

		if err != nil {
			log.Println("Failed to do mass get", err)
			continue
		}

		cars = append(cars, driverDetails...)
	}

	// at this point we have our cars alright, but it's in a stringified format
	// so the code below loop through our hits and converts them to array

	var parsedDrivers []interface{}

	for _, stringifiedCars := range cars {
		var parsedDriver interface{}

		str, isString := stringifiedCars.(string)
		if !isString {
			log.Println("Value retrieved from redis not string")
			continue
		}

		err := json.Unmarshal([]byte(str), &parsedDriver)
		if err != nil {
			log.Println("Failed to parse")
			continue
		}

		parsedDrivers = append(parsedDrivers, parsedDriver)
	}

	reponseValue := map[string]interface{}{
		"drivers": parsedDrivers,
	}

	return Response{
		Data:    reponseValue,
		Message: "Retrived closest drivers successfully",
	}
}

type DResponse struct {
	Response
	Data    []DriverWithTimeAndDistance
	Message string
}

func (c *DriverController) Dispatch(data TripRequestData) (DResponse, error) {

	response := c.FindClosestDrivers(data, 2)
	drivers := response.Data["drivers"].([]interface{})

	var transformedDrivers = []DriverInstance{}
	var lngLats = []string{}

	for _, value := range drivers {
		preDriver := value.(map[string]interface{})
		preCoordinates := preDriver["coordinates"].(map[string]interface{})

		driver := DriverInstance{
			Id: preDriver["id"].(string),
			Coordinates: GeoCoord{
				Latitude:  preCoordinates["latitude"].(float64),
				Longitude: preCoordinates["longitude"].(float64),
			},
			LastKnownIndex: preDriver["lastKnownIndex"].(string),
		}

		transformedDrivers = append(transformedDrivers, driver)
		stringifiedCoordinates := utils.StringifyLngLat(h3.GeoCoord{
			Latitude:  driver.Coordinates.Latitude,
			Longitude: driver.Coordinates.Longitude,
		})

		lngLats = append(lngLats, stringifiedCoordinates)

	}

	log.Println(lngLats)

	mapbox, err := scully.New(os.Getenv("ACCESS_TOKEN"))

	if err != nil {
		return DResponse{}, err
	}

	parsedValue := utils.ParseCoord(utils.IndexCoordinatesProps{
		Lat: data.Latitude,
		Lng: data.Longitude,
	})

	lngLats = append(lngLats, utils.StringifyLngLat(h3.GeoCoord{
		Latitude:  parsedValue.Lat,
		Longitude: parsedValue.Lng,
	}))

	points := strings.Join(lngLats[:], ";")

	mapbox.Matrix.SetDestinationIndex(strconv.Itoa(len(lngLats) - 1))

	mapboxResponse, err := mapbox.Matrix.GetMatrix(points)

	distances := mapboxResponse["distances"].([]interface{})

	durations := mapboxResponse["durations"].([]interface{})

	driverAndEtaData := []DriverWithTimeAndDistance{}

	for index, driver := range transformedDrivers {
		duration := durations[index].([]interface{})
		distance := distances[index].([]interface{})

		driverAndEtaData = append(driverAndEtaData, DriverWithTimeAndDistance{
			Driver: driver,
			DT: DistanceAndTime{
				Time:     duration[0].(float64),
				Distance: distance[0].(float64),
			},
		})

	}

	sort.Slice(driverAndEtaData, func(i, j int) bool {
		return driverAndEtaData[i].DT.Distance < driverAndEtaData[j].DT.Distance
	})

	if err != nil {
		return DResponse{}, err
	}

	return DResponse{
		Data:    driverAndEtaData,
		Message: "ETA Service",
	}, nil

}
