package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/jhiven/goat_data/overpass"
)

func writeOutput(elements *[]overpass.Element) {
	fmt.Println("Writing output")
	fo, err := os.Create("output.json")

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	buf := new(bytes.Buffer)

	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(elements); err != nil {
		panic(err)
	}

	if _, err := io.Copy(fo, buf); err != nil {
		panic(err)
	}
}

func main() {
	// jawa timur, jawa tengah, jawa barat, jakarta, yokyakarta, banten
	ids := []uint{3603438227, 3602388357, 3602388361, 3606362934, 3605616105, 3602388356}
	// Surabaya, Sidoarjo
	// ids := []uint{3608225862, 3609677345}
	overpassElements := make([]overpass.Element, 0)

	overpassCh := make(chan overpass.OverpassRes)
	postprocessCh := make(chan []overpass.Element)

	for _, id := range ids {
		go overpass.FetchOverpass(id, overpassCh)
		go overpass.PostProcessing(overpassCh, postprocessCh)
	}

	for range ids {
		overpassElements = append(overpassElements, <-postprocessCh...)
	}

	writeOutput(&overpassElements)
}
