package location

import (
	"encoding/json"
	"net/http"
	"time"
)

const locationsURL = "https://api-gw.criticalmaps.net/locations"

type Location struct {
	Device    string    `json:"device"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
}

func (l *Location) UnmarshalJSON(data []byte) error {
	type location struct {
		Device    string
		Latitude  int
		Longitude int
		Timestamp int64
	}

	loc := new(location)
	if err := json.Unmarshal(data, &loc); err != nil {
		return err
	}

	l.Device = loc.Device
	l.Latitude = float64(loc.Latitude) / 1000000
	l.Longitude = float64(loc.Longitude) / 1000000

	timezone, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	l.Timestamp = time.Unix(loc.Timestamp, 0).In(timezone)

	return nil
}

func FetchLocations() ([]Location, error) {
	resp, err := http.Get(locationsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic(resp.StatusCode)
	}

	locations := make([]Location, 0)
	if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
		return nil, err
	}

	return locations, nil
}
