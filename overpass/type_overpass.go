package overpass

import "github.com/jhiven/goat_data/haversine"

type Osm3sData struct {
	TimestampOsm  string `json:"timestamp_osm_base"`
	TimestampArea string `json:"timestamp_areas_base"`
}

type Tag struct {
	Name        *string `json:"name,omitempty"`
	City        *string `json:"addr:city,omitempty"`
	Address     *string `json:"addr:address,omitempty"`
	Postcode    *string `json:"addr:postcode,omitempty"`
	HouseNumber *string `json:"addr:house_number,omitempty"`
	Street      *string `json:"addr:street,omitempty"`
	Brand       *string `json:"brand,omitempty"`
	Shop        *string `json:"shop,omitempty"`
	Amenity     *string `json:"amenity,omitempty"`
	Operator    *string `json:"operator,omitempty"`
	Building    *string `json:"building,omitempty"`
}

type Center struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Element struct {
	ElementId uint     `json:"id"`
	Type      string   `json:"type"`
	Lat       *float64 `json:"lat,omitempty"`
	Lon       *float64 `json:"lon,omitempty"`
	Center    *Center  `json:"center,omitempty"`
	Nodes     *[]uint  `json:"nodes,omitempty"`
	Tags      Tag      `json:"tags"`
}

type OverpassRes struct {
	Version   float32   `json:"version"`
	Generator string    `json:"generator"`
	Osm3s     Osm3sData `json:"osm3s"`
	Elements  []Element `json:"elements"`
}

func (e *Element) GetLanLon() haversine.LatLon {
	if e.Lat != nil && e.Lon != nil {
		return haversine.LatLon{Lat: *e.Lat, Lon: *e.Lon}
	}

	return haversine.LatLon{Lat: e.Center.Lat, Lon: e.Center.Lon}
}
