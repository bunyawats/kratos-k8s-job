package common

import (
	"fmt"
	"runtime/metrics"
	"strconv"
)

func GetGoRuntimeMetrics() (map[string]interface{}, error) {

	currentMatrix := make(map[string]interface{})

	// Get descriptions for all supported metrics.
	descs := metrics.All()

	// Create a sample for each metric.
	samples := make([]metrics.Sample, len(descs))
	for i := range samples {
		samples[i].Name = descs[i].Name
	}

	// Sample the metrics. Re-use the samples slice if you can!
	metrics.Read(samples)

	// Iterate over all results.
	for _, sample := range samples {

		name, value := sample.Name, sample.Value

		switch value.Kind() {
		case metrics.KindUint64:
			currentMatrix[name] = value.Uint64()
		case metrics.KindFloat64:
			currentMatrix[name] = value.Float64()
		case metrics.KindFloat64Histogram:
			medianBk := medianBucket(value.Float64Histogram())
			currentMatrix[name] = strconv.FormatFloat(medianBk, 'f', -1, 64)
		case metrics.KindBad:
			fmt.Println("bug in runtime/metrics package!")
		default:
			fmt.Printf("%s: unexpected metric Kind: %v\n", name, value.Kind())
		}
		fmt.Printf("%s: %v\n", name, currentMatrix[name])
	}

	return currentMatrix, nil
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
