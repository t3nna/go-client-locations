package util

import (
	"go-clinet-locations/shared/types"
	"math"
	"testing"
)

func TestCalculateDistance(t *testing.T) {
	tests := []struct {
		name         string
		coord1       *types.Coordinate
		coord2       *types.Coordinate
		expectedDist float64
		allowedError float64
	}{
		{
			name: "same coordinates",
			coord1: &types.Coordinate{
				Latitude:  51.11822470712269,
				Longitude: 16.990711729269563,
			},
			coord2: &types.Coordinate{
				Latitude:  51.11822470712269,
				Longitude: 16.990711729269563,
			},
			expectedDist: 0.0,
			allowedError: 0.001,
		},
		{
			name: "short distance - same city",
			coord1: &types.Coordinate{
				Latitude:  51.11822470712269,
				Longitude: 16.990711729269563,
			},
			coord2: &types.Coordinate{
				Latitude:  51.11956092410769,
				Longitude: 17.05696305051491,
			},
			expectedDist: 4.5, // Approximately 4.5 km
			allowedError: 0.5,
		},
		{
			name: "medium distance - between cities",
			coord1: &types.Coordinate{
				Latitude:  51.11822470712269, // Wroclaw
				Longitude: 16.990711729269563,
			},
			coord2: &types.Coordinate{
				Latitude:  52.23553956649786, // Warsaw
				Longitude: 20.984595191389918,
			},
			expectedDist: 300.0, // Approximately 300 km
			allowedError: 50.0,
		},
		{
			name: "long distance - different countries",
			coord1: &types.Coordinate{
				Latitude:  51.11822470712269, // Wroclaw, Poland
				Longitude: 16.990711729269563,
			},
			coord2: &types.Coordinate{
				Latitude:  50.53401932980686, // Kyiv, Ukraine
				Longitude: 31.178889172903055,
			},
			expectedDist: 1000.0, // Approximately 1000 km
			allowedError: 100.0,
		},
		{
			name: "antipodal points",
			coord1: &types.Coordinate{
				Latitude:  0.0,
				Longitude: 0.0,
			},
			coord2: &types.Coordinate{
				Latitude:  0.0,
				Longitude: 180.0,
			},
			expectedDist: 20015.0, // Half the Earth's circumference
			allowedError: 100.0,
		},
		{
			name: "equator distance",
			coord1: &types.Coordinate{
				Latitude:  0.0,
				Longitude: 0.0,
			},
			coord2: &types.Coordinate{
				Latitude:  0.0,
				Longitude: 1.0,
			},
			expectedDist: 111.0, // Approximately 111 km per degree at equator
			allowedError: 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance := CalculateDistance(tt.coord1, tt.coord2)

			if math.Abs(distance-tt.expectedDist) > tt.allowedError {
				t.Errorf("expected distance around %.2f km, got %.2f km (allowed error: %.2f)",
					tt.expectedDist, distance, tt.allowedError)
			}
		})
	}
}

func TestCalculateDistanceCommutativity(t *testing.T) {
	coord1 := &types.Coordinate{
		Latitude:  51.11822470712269,
		Longitude: 16.990711729269563,
	}
	coord2 := &types.Coordinate{
		Latitude:  52.23553956649786,
		Longitude: 20.984595191389918,
	}

	distance1 := CalculateDistance(coord1, coord2)
	distance2 := CalculateDistance(coord2, coord1)

	if math.Abs(distance1-distance2) > 0.001 {
		t.Errorf("distance calculation should be commutative, got %.6f and %.6f", distance1, distance2)
	}
}

func TestCalculateDistanceWithZeroCoordinates(t *testing.T) {
	coord1 := &types.Coordinate{
		Latitude:  0.0,
		Longitude: 0.0,
	}
	coord2 := &types.Coordinate{
		Latitude:  0.0,
		Longitude: 0.0,
	}

	distance := CalculateDistance(coord1, coord2)
	if distance != 0.0 {
		t.Errorf("distance between same coordinates should be 0, got %.6f", distance)
	}
}
