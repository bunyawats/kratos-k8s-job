package main

import (
	"context"
	"fmt"
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"os"
	"time"
)

func main() {
	// Create client
	url := "https://us-east-1-1.aws.cloud2.influxdata.com"
	//token := "tOcSMS46PCpylMuTbvDNMt6iZs-7YlsAbsjXHXPIsNuKMFHgWgLNE9Zp0ukNnfvtkAhVXJqc_tOApztwYdy6mQ=="
	token := os.Getenv("INFLUXDB_API_KEY")
	// Create a new client using an InfluxDB server base URL and an authentication token
	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:  url,
		Token: token,
	})

	if err != nil {
		panic(err)
	}
	// Close client at the end and escalate error if present
	defer func(client *influxdb3.Client) {
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}(client)

	database := "k8s-job"

	data := map[string]map[string]interface{}{
		"point1": {
			"location": "Klamath",
			"species":  "bees",
			"count":    23,
		},
		"point2": {
			"location": "Portland",
			"species":  "ants",
			"count":    30,
		},
		"point3": {
			"location": "Klamath",
			"species":  "bees",
			"count":    28,
		},
		"point4": {
			"location": "Portland",
			"species":  "ants",
			"count":    32,
		},
		"point5": {
			"location": "Klamath",
			"species":  "bees",
			"count":    29,
		},
		"point6": {
			"location": "Portland",
			"species":  "ants",
			"count":    40,
		},
	}
	fmt.Print("data \n", data)

	//Write data
	options := influxdb3.WriteOptions{
		Database: database,
	}
	for key := range data {
		point := influxdb3.NewPointWithMeasurement("census").
			SetTag("location", data[key]["location"].(string)).
			SetField(data[key]["species"].(string), data[key]["count"])

		if err := client.WritePointsWithOptions(context.Background(), &options, point); err != nil {
			panic(err)
		}

		time.Sleep(1 * time.Second) // separate points by 1 second
	}

	// Execute query
	query := `SELECT *
          FROM 'census'
          WHERE time >= now() - interval '12 hour'
            AND ('bees' IS NOT NULL OR 'ants' IS NOT NULL)`

	queryOptions := influxdb3.QueryOptions{
		Database: database,
	}
	iterator, err := client.QueryWithOptions(context.Background(), &queryOptions, query)

	if err != nil {
		panic(err)
	}

	for iterator.Next() {
		value := iterator.Value()

		location := value["location"]
		ants := value["ants"]
		bees := value["bees"]
		fmt.Printf("in %s are %d ants and %d bees\n", location, ants, bees)
	}

}
