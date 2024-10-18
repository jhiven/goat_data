package overpass

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/jhiven/goat_data/internal/haversine"
)

func TestOutputEachDistance(t *testing.T) {
	var elements []Element
	totalFail := 0
	res, err := os.ReadFile("../output.json")
	if err != nil {
		fmt.Println("error open file")
		panic(err)
	}
	if err := json.Unmarshal(res, &elements); err != nil {
		fmt.Println("error unmarshaling")
		panic(err)
	}

	for i, element := range elements {
		for j, item := range elements {
			if i == j {
				continue
			}
			got := haversine.DistanceInM(element.GetLanLon(), item.GetLanLon())

			if got < 20 {
				totalFail++
				t.Errorf(
					"fail: distance [%v] and [%v] is %v",
					*element.Tags.Name,
					*item.Tags.Name,
					got,
				)
			}
		}
	}

	if totalFail > 0 {
		t.Errorf("total fail %v", totalFail)
	}
}

func TestFilterOutput(t *testing.T) {
	overpassCh := make(chan OverpassRes)
	postprocessCh := make(chan []Element)

	elements := GetElements()
	go func() {
		overpassCh <- OverpassRes{Elements: elements}
	}()
	go PostProcessing(overpassCh, postprocessCh)

	postProcessLen := len(<-postprocessCh)
	overpassLen := len(elements)

	if overpassLen != postProcessLen {
		t.Errorf("fail: expected total item: %v, got: %v", overpassLen, postProcessLen)
	}
}
