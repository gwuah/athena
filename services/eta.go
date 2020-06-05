package services

import (
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/electra-systems/athena/storage"

	"github.com/electra-systems/athena/utils"
	"github.com/gwuah/scully"
	"github.com/uber/h3-go"
)

type ETA struct {
	DB           map[string]storage.Redis
	ResourceName string
	IndexName    string
}

func NewETA(db map[string]storage.Redis, ResourceName string, IndexName string) *ETA {
	return &ETA{DB: db, ResourceName: ResourceName, IndexName: IndexName}
}

func transformRawEntityData(rawData []interface{}) ([]string, []Entity) {

	var transformedEntities = []Entity{}
	var lngLats = []string{}

	for _, value := range rawData {
		preDriver := value.(map[string]interface{})
		preCoordinates := preDriver["coordinates"].(map[string]interface{})

		entity := Entity{
			Id: preDriver["id"].(string),
			Coordinates: GeoCoord{
				Lat: preCoordinates["lat"].(float64),
				Lng: preCoordinates["lng"].(float64),
			},
			LastKnownIndex: preDriver["lastKnownIndex"].(string),
		}

		transformedEntities = append(transformedEntities, entity)
		stringifiedCoordinates := utils.StringifyLngLat(h3.GeoCoord{
			Latitude:  entity.Coordinates.Lat,
			Longitude: entity.Coordinates.Lng,
		})

		lngLats = append(lngLats, stringifiedCoordinates)

	}

	return lngLats, transformedEntities

}

func (e *ETA) GetEta(data Payload, entities []interface{}, sortyBy string) ([]EntityWithETA, error) {

	if len(entities) == 0 {
		return []EntityWithETA{}, nil
	}

	lngLats, transformedEntities := transformRawEntityData(entities)

	parsedValue := utils.ParseCoord(utils.IndexCoordinatesProps{
		Lat: data.Lat,
		Lng: data.Lng,
	})

	// attach origin lng/lat to request
	lngLats = append(lngLats, utils.StringifyLngLat(h3.GeoCoord{
		Latitude:  parsedValue.Lat,
		Longitude: parsedValue.Lng,
	}))

	points := strings.Join(lngLats[:], ";")

	mapbox, err := scully.New(os.Getenv("ACCESS_TOKEN"))

	if err != nil {
		return nil, err
	}

	mapbox.Matrix.SetDestinationIndex(strconv.Itoa(len(lngLats) - 1))

	mapboxResponse, err := mapbox.Matrix.GetMatrix(points)

	if err != nil {
		return nil, err
	}

	distances := mapboxResponse["distances"].([]interface{})

	durations := mapboxResponse["durations"].([]interface{})

	entitiesWithEtaAttached := []EntityWithETA{}

	for index, entity := range transformedEntities {
		duration := durations[index].([]interface{})
		distance := distances[index].([]interface{})

		entitiesWithEtaAttached = append(entitiesWithEtaAttached, EntityWithETA{
			Entity: entity,
			DT: DistanceAndTime{
				Time:     duration[0].(float64),
				Distance: distance[0].(float64),
			},
		})

	}

	if sortyBy == "distance" {
		sort.Slice(entitiesWithEtaAttached, func(i, j int) bool {
			return entitiesWithEtaAttached[i].DT.Distance < entitiesWithEtaAttached[j].DT.Distance
		})
	} else {
		sort.Slice(entitiesWithEtaAttached, func(i, j int) bool {
			return entitiesWithEtaAttached[i].DT.Time < entitiesWithEtaAttached[j].DT.Time
		})
	}

	return entitiesWithEtaAttached, nil
}
