# Athena

Athena is a an API service that provides 4 main services, geo-indexing, dispatching, proximity-searching, ETA. <br/>
We leverage mapbox for features such as distance-matrix to enable us sort hits in ascending order. <br/>
It's built to work as a core service that order services will integrate with. <br/>

# Deps
- github.com/gwuah/scully
- github.com/joho/godotenv
- github.com/gin-gonic/gin
- github.com/uber/h3-go
- github.com/go-redis/redis

# Status

- [x] geo-indexing
- [x] searching
- [x] ETA
- [] dispatch (Currently working on this)

# Core Modules

For the sake of reusablilty, I've placed the core modules them in a services folder where they can be imported and used.. <br/>

Example usage of ETA service

- `services.NewETA(databaseMap, resourceName, indexName)`<br/>
  Where,
- `indexName` - Index in which a search query may be made
- `resourceName` - Name of a resource

Eg. Drivers are resources. But the index in which we'll look for drivers is called the "supplyIndex". So properties(id, lng/lat) etc of the driver is stored in the resource database and they are geo-indexed in a seperate database, which we are calling the supplyIndex.

# How to run

First Way

- `go run server.go`

Second Way

- `go build server.go` and then you run `./server`

# Usage

To index driver location data,

`POST /index-driver-location`

```
{
    "id": "8",
    "lat": "5.678787197821624",
    "lng": "-0.25505293160676956"
}
```

Response

```
{
    "data": {
        "coordinates": {
            "lat": 5.678787197821624,
            "lng": -0.25505293160676956
        },
        "current_index": "614555982713847807",
        "previous_index": "614555982713847807"
        "id": "8",
    },
    "message": "Successfully Indexed Coordinates"
}
```

To find closest drivers,

`POST /closest-drivers`

```
{
    "id": "3",
    "lat": "5.678787197821624",
    "lng": "-0.25505293160676956"
}
```

Response

```
{
    "data": [
        {
            "entity": {
                "id": "3",
                "lastKnownIndex": "614555982713847807",
                "coordinates": {
                    "lat": 5.678787197821624,
                    "lng": -0.25505293160676956
                }
            },
            "dt": {
                "time": 0,
                "distance": 0
            }
        },
        {
            "entity": {
                "id": "8",
                "lastKnownIndex": "614555982741110783",
                "coordinates": {
                    "lat": 5.676310318306305,
                    "lng": -0.24685610085725784
                }
            },
            "dt": {
                "time": 302.9,
                "distance": 1253.1
            }
        },
        {
            "entity": {
                "id": "1",
                "lastKnownIndex": "614555982741110783",
                "coordinates": {
                    "lat": 5.676310318306305,
                    "lng": -0.24685610085725784
                }
            },
            "dt": {
                "time": 302.9,
                "distance": 1253.1
            }
        },
        {
            "entity": {
                "id": "2",
                "lastKnownIndex": "614555982722236415",
                "coordinates": {
                    "lat": 5.684125264018471,
                    "lng": -0.24913061410188675
                }
            },
            "dt": {
                "time": 264.8,
                "distance": 1312.9
            }
        },
        {
            "entity": {
                "id": "6",
                "lastKnownIndex": "614555983238135807",
                "coordinates": {
                    "lat": 5.680324565958827,
                    "lng": -0.2416633442044258
                }
            },
            "dt": {
                "time": 466.8,
                "distance": 2468.4
            }
        },
        {
            "entity": {
                "id": "0",
                "lastKnownIndex": "614555982736916479",
                "coordinates": {
                    "lat": 5.674986464545531,
                    "lng": -0.23887384682893753
                }
            },
            "dt": {
                "time": 633.7,
                "distance": 3265.3
            }
        }
    ],
    "message": "Closest Drivers Retrieved"
}
```
