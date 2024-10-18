package haversine

import (
	"testing"
)

type testCoordinates struct {
	p            LatLon
	q            LatLon
	realDistance float64
}

var tests = []testCoordinates{
	{
		LatLon{22.55, 43.12},  // Rio de Janeiro, Brazil
		LatLon{13.45, 100.28}, // Bangkok, Thailand
		6094.544408786774,
	},
	{
		LatLon{20.10, 57.30}, // Port Louis, Mauritius
		LatLon{0.57, 100.21}, // Padang, Indonesia
		5145.525771394785,
	},
	{
		LatLon{51.45, 1.15},  // Oxford, United Kingdom
		LatLon{41.54, 12.27}, // Vatican, City Vatican City
		1389.1793118293067,
	},
	{
		LatLon{22.34, 17.05}, // Windhoek, Namibia
		LatLon{51.56, 4.29},  // Rotterdam, Netherlands
		3429.89310043882,
	},
	{
		LatLon{63.24, 56.59}, // Esperanza, Argentina
		LatLon{8.50, 13.14},  // Luanda, Angola
		6996.18595539861,
	},
	{
		LatLon{90.00, 0.00}, // North/South Poles
		LatLon{48.51, 2.21}, // Paris,  France
		4613.477506482742,
	},
	{
		LatLon{45.04, 7.42},  // Turin, Italy
		LatLon{3.09, 101.42}, // Kuala Lumpur, Malaysia
		10078.111954385415,
	},
}

func TestDistance(t *testing.T) {
	for _, input := range tests {
		km := DistanceInKm(input.p, input.q)

		if input.realDistance != km {
			t.Errorf("fail: want %v %v -> %v got %v",
				input.p,
				input.q,
				input.realDistance,
				km,
			)
		}
	}

}
