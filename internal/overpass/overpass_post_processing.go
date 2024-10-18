package overpass

import (
	"fmt"

	"github.com/jhiven/goat_data/internal/haversine"
)

const minDistance float64 = 20

func PostProcessing(overpassCh <-chan OverpassRes, postprocessCh chan<- []Element) {
	elements := (<-overpassCh).Elements
	results := make([]Element, 0)
	fmt.Println("Running post processing")

	for _, element := range elements {
		latlon1 := element.GetLanLon()
		isValid := true

		for _, result := range results {
			latlon2 := result.GetLanLon()
			distance := haversine.DistanceInM(latlon1, latlon2)

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
			results = append(results, element)
		}
	}

	fmt.Println("Post processing DONE!")
	fmt.Printf("before: %v, after: %v\n", len(elements), len(results))

	postprocessCh <- results
}

func RemoveDuplicate(sliceList []Element) []Element {
	fmt.Println("Removing duplicates")
	allKeys := make(map[uint]struct{})
	list := make([]Element, 0)
	for _, item := range sliceList {
		if _, value := allKeys[item.ElementId]; !value {
			allKeys[item.ElementId] = struct{}{}
			list = append(list, item)
		} else {
			fmt.Printf(
				"Duplicate id: %v, name: %v\n",
				item.ElementId,
				*item.Tags.Name,
			)
		}
	}
	return list
}
