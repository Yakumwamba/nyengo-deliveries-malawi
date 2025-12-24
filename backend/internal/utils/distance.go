package utils

import "math"

// Haversine calculates distance between two coordinates in kilometers
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // km

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// RoadDistance estimates road distance (straight line * factor)
func RoadDistance(lat1, lon1, lat2, lon2 float64) float64 {
	return Haversine(lat1, lon1, lat2, lon2) * 1.3
}

// EstimateDuration estimates travel time in minutes
func EstimateDuration(distanceKm float64, avgSpeedKmh float64) int {
	if avgSpeedKmh <= 0 {
		avgSpeedKmh = 30 // default urban speed
	}
	minutes := int(distanceKm / avgSpeedKmh * 60)
	if minutes < 10 {
		return 10
	}
	return minutes
}
