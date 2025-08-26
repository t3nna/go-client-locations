package util

import (
	"go-clinet-locations/shared/types"
	"math"
)

// calculateDistance computes the distance between two coordinates using the Haversine formula.

func CalculateDistance(coord1, coord2 *types.Coordinate) float64 {
	const earthRadius = 6371 // Earth's radius in kilometers

	lat1, lon1 := degreesToRadians(coord1.Latitude), degreesToRadians(coord1.Longitude)
	lat2, lon2 := degreesToRadians(coord2.Latitude), degreesToRadians(coord2.Longitude)

	deltaLat := lat2 - lat1
	deltaLon := lon2 - lon1

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// degreesToRadians converts degrees to radians.
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
