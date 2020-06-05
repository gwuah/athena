package services

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/electra-systems/athena/storage"
	"github.com/electra-systems/athena/utils"
	"github.com/go-redis/redis"
)

type GeoIndex struct {
	DB           map[string]storage.Redis
	ResourceName string
	IndexName    string
}

type Response struct {
	Err  error
	Data map[string]interface{}
}

type GeoCoord struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func NewGeoIndex(db map[string]storage.Redis, ResourceName string, IndexName string) *GeoIndex {
	return &GeoIndex{DB: db, ResourceName: ResourceName, IndexName: IndexName}
}

func (g *GeoIndex) Index(data Payload) (map[string]interface{}, error) {

	if g.DB[g.ResourceName] == nil {
		err := errors.New("Resource [" + g.ResourceName + "] doesn't exist!")
		log.Println(err)
		return nil, err
	}

	indexingResponse := utils.IndexCoordinates(utils.IndexCoordinatesProps{
		Lat: data.Lat,
		Lng: data.Lng,
	})

	stringifiedIndex := utils.H3IndexToString(indexingResponse.Index)

	response, err := g.DB[g.ResourceName].Get(data.Id)

	if err != redis.Nil && err != nil {
		log.Println(err)
		return nil, err
	}

	var oldEntityInstance Entity

	err = json.Unmarshal([]byte(response), &oldEntityInstance)

	if stringifiedIndex == oldEntityInstance.LastKnownIndex {
		err := errors.New("Entity hasn't changed their location")
		log.Println(err)

		return map[string]interface{}{
			"id":             oldEntityInstance.Id,
			"previous_index": oldEntityInstance.LastKnownIndex,
			"current_index":  oldEntityInstance.LastKnownIndex,
			"coordinates": map[string]interface{}{
				"lat": indexingResponse.Lat,
				"lng": indexingResponse.Lng,
			},
		}, nil
	}

	newEntityInstance := Entity{
		Id:             data.Id,
		LastKnownIndex: stringifiedIndex,
		Coordinates: GeoCoord{
			Lat: indexingResponse.Lat,
			Lng: indexingResponse.Lng,
		},
	}

	marshalledValue, err := json.Marshal(newEntityInstance)

	_, err = g.DB[g.ResourceName].Set(newEntityInstance.Id, marshalledValue)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	_, err = g.DB[g.IndexName].RemoveFromList(oldEntityInstance.LastKnownIndex, newEntityInstance.Id)

	if err != nil {
		err := errors.New("Failed to update old index")
		log.Println(err)
		return nil, err
	}

	_, err = g.DB[g.IndexName].InsertIntoList(stringifiedIndex, data.Id)

	if err != nil {
		err := errors.New("Failed to update new index")
		log.Println(err)
		return nil, err
	}

	reponseValue := map[string]interface{}{
		"id":             newEntityInstance.Id,
		"previous_index": oldEntityInstance.LastKnownIndex,
		"current_index":  newEntityInstance.LastKnownIndex,
		"coordinates": map[string]interface{}{
			"lat": indexingResponse.Lat,
			"lng": indexingResponse.Lng,
		},
	}

	return reponseValue, nil

}
