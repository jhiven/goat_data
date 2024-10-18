# GoAT Data Processing

This repository contains geospatial data processing tools for the product [GoAT](https://goat.agileteknik.com/), implemented in Go. It processes data from **OpenStreetMap** and **Overpass Turbo**, ensuring quality by removing duplicate entries and maintaining a minimum distance of 20 meters between points.

## Features

- Fetch and clean geospatial data from OpenStreetMap.
- Remove duplicate data points.
- Ensure each point has at least 20m distance from others.

## Usage

1. Install dependencies:
   ```bash
   go mod tidy
   ```
2. Run the main program:
   ```bash
   go run main.go
   ```

## Data Sources

- [OpenStreetMap](https://www.openstreetmap.org/)
- [Overpass Turbo](https://overpass-turbo.eu/)
