package overpass

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
)

func FetchOverpass(id uint, ch chan<- OverpassRes) {

	overpassUrl := "https://overpass-api.de/api/interpreter"

	payload := fmt.Sprintf(`
[out:json][timeout:10000];
area(id:%v)->.searchArea;
(
node["building"]["name"](area.searchArea);
way["building"]["name"](area.searchArea);
node["amenity"]["name"](area.searchArea);
way["amenity"]["name"](area.searchArea);
node["shop"]["name"](area.searchArea);
way["shop"]["name"](area.searchArea);
node["office"]["name"](area.searchArea);
way["office"]["name"](area.searchArea);
node["place"]["name"](area.searchArea);
way["place"]["name"](area.searchArea);
node["public_transport"]["name"](area.searchArea);
way["public_transport"]["name"](area.searchArea);
);
out center;
`, id)

	fmt.Printf("Fetching area id: %v\n", id)

	res, err := http.PostForm(overpassUrl, url.Values{"data": {payload}})
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Got data with area id: %v!\n", id)

	var overpass OverpassRes

	fmt.Printf("Decoding area id: %v\n", id)
	if err := json.Unmarshal(body, &overpass); err != nil {
		fmt.Println("error unmarshaling")
		panic(err)
	}
	fmt.Printf("Area id %v DONE!\n", id)

	ch <- overpass
}

// Testing with file purpose
func GetElements() []Element {
	var elements []Element

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	res, err := os.ReadFile(fmt.Sprintf("%v/../output.json", basepath))

	if err != nil {
		fmt.Println("error open file")
		panic(err)
	}
	if err := json.Unmarshal(res, &elements); err != nil {
		fmt.Println("error unmarshaling")
		panic(err)
	}

	return elements
}
