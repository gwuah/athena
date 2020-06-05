package services

import (
	"encoding/json"
	"log"

	"github.com/electra-systems/athena/storage"
	"github.com/electra-systems/athena/utils"
	"github.com/uber/h3-go"
)

type Search struct {
	DB           map[string]storage.Redis
	ResourceName string
	IndexName    string
}

func NewSearch(db map[string]storage.Redis, ResourceName string, IndexName string) *Search {
	return &Search{DB: db, ResourceName: ResourceName, IndexName: IndexName}
}

func (s *Search) Closest(data Payload, neighbours int) ([]interface{}, error) {
	var indexResponse = utils.IndexCoordinates(utils.IndexCoordinatesProps{
		Lat: data.Lat,
		Lng: data.Lng,
	})

	rings := h3.KRing(indexResponse.Index, neighbours)

	results := []interface{}{}

	for _, value := range rings {
		entityIds, err := s.DB[s.IndexName].All(utils.H3IndexToString(value))

		if err != nil {
			log.Println(err)
			continue
		}

		if len(entityIds) == 0 {
			// if we pass an empty array to Mget, it throws an error
			// so we just exit the current iteration and move on
			continue
		}

		entityDetails, err := s.DB[s.ResourceName].MGet(entityIds)

		if err != nil {
			log.Println("Failed to do mass get", err)
			continue
		}

		results = append(results, entityDetails...)
	}

	if len(results) == 0 {
		return results, nil
	}

	// at this point we have our results alright, but it's in a stringified format
	// so the code below loop through our hits and converts them to array

	var parsedEntities []interface{}

	for _, stringifiedCars := range results {
		var parsedEntity interface{}

		str, isString := stringifiedCars.(string)
		if !isString {
			log.Println("Value retrieved from redis not string")
			continue
		}

		err := json.Unmarshal([]byte(str), &parsedEntity)
		if err != nil {
			log.Println("Failed to parse")
			continue
		}

		parsedEntities = append(parsedEntities, parsedEntity)
	}

	return parsedEntities, nil
}
