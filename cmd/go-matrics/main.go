package main

import (
	"context"
	"fmt"
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"runtime/metrics"
	"strconv"
	"time"
)

func main() {

	url := "https://us-east-1-1.aws.cloud2.influxdata.com"
	token := "tOcSMS46PCpylMuTbvDNMt6iZs-7YlsAbsjXHXPIsNuKMFHgWgLNE9Zp0ukNnfvtkAhVXJqc_tOApztwYdy6mQ=="

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

	options := influxdb3.WriteOptions{
		Database: database,
	}

	// Get descriptions for all supported metrics.
	descs := metrics.All()

	// Create a sample for each metric.
	samples := make([]metrics.Sample, len(descs))
	for i := range samples {
		samples[i].Name = descs[i].Name
	}

	for i := 0; i < 10; i++ {

		// Sample the metrics. Re-use the samples slice if you can!
		metrics.Read(samples)

		point := influxdb3.NewPointWithMeasurement("metrics")

		// Iterate over all results.
		for _, sample := range samples {

			// Pull out the name and value.
			name, value := sample.Name, sample.Value

			// Handle each sample.
			switch value.Kind() {
			case metrics.KindUint64:
				fmt.Printf("KindUint64: %s: %d\n", name, value.Uint64())
				point.SetField(name, value.Uint64())
			case metrics.KindFloat64:
				fmt.Printf("KindFloat64: %s: %f\n", name, value.Float64())
				point.SetField(name, value.Float64())
			case metrics.KindFloat64Histogram:
				// The histogram may be quite large, so let's just pull out
				// a crude estimate for the median for the sake of this example.

				//fmt.Println("value.Float64Histogram: ", value.Float64Histogram())
				medianBk := medianBucket(value.Float64Histogram())
				fmt.Printf("KindFloat64Histogram: %s: %f\n", name, medianBk)
				point.SetTag(name, strconv.FormatFloat(medianBk, 'f', -1, 64))
			case metrics.KindBad:
				// This should never happen because all metrics are supported
				// by construction.
				//panic("bug in runtime/metrics package!")
				fmt.Println("bug in runtime/metrics package!")
			default:
				// This may happen as new metrics get added.
				//
				// The safest thing to do here is to simply log it somewhere
				// as something to look into, but ignore it for now.
				// In the worst case, you might temporarily miss out on a new metric.
				fmt.Printf("%s: unexpected metric Kind: %v\n", name, value.Kind())
			}
		}
		if err := client.WritePointsWithOptions(context.Background(), &options, point); err != nil {
			fmt.Println("error while writing point to InfluxDB", err.Error())
		}
		fmt.Println("index: ", i, "\n")

		time.Sleep(time.Second)
	}

}

func medianBucket(h *metrics.Float64Histogram) float64 {
	total := uint64(0)
	for _, count := range h.Counts {
		total += count
	}
	thresh := total / 2
	total = 0
	for i, count := range h.Counts {
		total += count
		if total >= thresh {
			return h.Buckets[i]
		}
	}
	panic("should not happen")
}
