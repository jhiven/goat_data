package overpass

import (
	"fmt"

	"github.com/jhiven/goat_data/haversine"
)

const minDistance float64 = 20

func PostProcessing(overpassCh <-chan OverpassRes, postprocessCh chan<- []Element) {
	results := make([]Element, 0)
	elements := (<-overpassCh).Elements
	fmt.Println("Running post processing")

	for i, element := range elements {
		var latlon2 haversine.LatLon
		latlon1 := element.GetLanLon()
		isValid := true

		for j, result := range results {
			if i == 0 && j == 0 {
				continue
			}

			latlon2 = result.GetLanLon()
			distance := haversine.DistanceInM(latlon1, latlon2)

			// fmt.Printf("element: %v, item: %v, distance: %v\n", *e.Tags.Name, *item.Tags.Name, m)
			if distance < minDistance {
				fmt.Printf(
					"[%v] is close to [%v] with distance: %v, REMOVING %v... \n",
					*element.Tags.Name,
					*result.Tags.Name,
					distance,
					*element.Tags.Name,
				)
				isValid = false
				break
			}

			if element.ElementId == result.ElementId {
				fmt.Printf(
					"Duplicate id with id: %v, name: %v\n",
					element.ElementId,
					*element.Tags.Name,
				)
				isValid = false
				break
			}
		}

		if isValid {
			// fmt.Printf("Appending element: %v\n", *e.Tags.Name)
			results = append(results, element)
		}
	}

	fmt.Println("Post processing DONE!")
	fmt.Printf("before: %v, after: %v\n", len(elements), len(results))

	postprocessCh <- results
}
