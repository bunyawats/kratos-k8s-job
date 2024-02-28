package common

import (
	"fmt"
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"runtime/metrics"
)

func GetGoRuntimeMetrics() ([]*influxdb3.Point, error) {
	//currentMatrix := make(map[string]interface{})

	// Get descriptions for all supported metrics.
	descs := metrics.All()

	// Create a sample for each metric.
	samples := make([]metrics.Sample, len(descs))
	for i := range samples {
		samples[i].Name = descs[i].Name
	}

	pts := make([]*influxdb3.Point, 0)

	// Sample the metrics. Re-use the samples slice if you can!
	metrics.Read(samples)

	for _, sample := range samples {

		// Pull out the name and value.
		name, value := sample.Name, sample.Value

		// Handle each sample.
		switch value.Kind() {
		case metrics.KindUint64:
			fmt.Printf("KindUint64\n: %s\n: %d\n", name, value.Uint64())
			point := influxdb3.NewPointWithMeasurement("metrics-simple")
			point.SetField("KindUint64", value.Uint64())
			point.SetTag("Name", name)
			pts = append(pts, point)
		case metrics.KindFloat64:
			fmt.Printf("KindFloat64\n: %s\n: %f\n", name, value.Float64())
			point := influxdb3.NewPointWithMeasurement("metrics-simple")
			point.SetField("KindFloat64", value.Float64())
			point.SetTag("Name", name)
			pts = append(pts, point)
		case metrics.KindFloat64Histogram:
			fmt.Printf("\nKindFloat64Histogram\n: %s\n: %v\n", name, value.Float64Histogram())
			point := influxdb3.NewPointWithMeasurement("metrics-simple")

			const maxBucketLen = 70

			startIndex := 0
			bucketLen := len(value.Float64Histogram().Buckets) - 1
			if bucketLen > maxBucketLen {
				midIndex := bucketLen / 2
				startIndex = midIndex - (maxBucketLen / 2)
				bucketLen = maxBucketLen
			}

			fmt.Println("bucketLen: ", bucketLen)
			fmt.Println("startIndex: ", startIndex)
			for i := 0; i < bucketLen; i++ {
				index := i - startIndex
				if index >= 0 && index < bucketLen {

					countValue := value.Float64Histogram().Counts[i]
					bucketValue := value.Float64Histogram().Buckets[i]
					if bucketValue < 0 {
						bucketValue = 0.0
					}

					point.SetField(fmt.Sprintf("count[%v]", index), countValue)
					point.SetField(fmt.Sprintf("bucket[%v]", index), bucketValue)
				}

			}
			point.SetTag("Name", name)
			pts = append(pts, point)

		case metrics.KindBad:
			fmt.Println("bug in runtime/metrics package!")
		default:
			fmt.Printf("%s: unexpected metric Kind: %v\n", name, value.Kind())
		}
	}

	return pts, nil
}
