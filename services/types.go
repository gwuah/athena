package services

type Payload struct {
	Id  string `json:"id"`
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type Entity struct {
	Id             string   `json:"id"`
	LastKnownIndex string   `json:"lastKnownIndex"`
	Coordinates    GeoCoord `json:"coordinates"`
}

type DistanceAndTime struct {
	Time     float64 `json:"time"`
	Distance float64 `json:"distance"`
}

type EntityWithETA struct {
	Entity Entity          `json:"entity"`
	DT     DistanceAndTime `json:"dt"`
}
