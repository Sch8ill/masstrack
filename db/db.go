package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/sch8ill/masstrack/location"

	_ "github.com/lib/pq"
)

type DB struct {
	db *sql.DB
	mu sync.Mutex
}

func NewDB(url string) (*DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("SET TIME ZONE 'Europe/Berlin'")
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func (db *DB) NewLocation(l location.Location) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, err := db.db.Exec("INSERT INTO locations (device, latitude, longitude, timestamp) SELECT $1::text, $2::double precision, $3::double precision, $4::timestamp WHERE NOT EXISTS ( SELECT 1 FROM locations WHERE device = $1::text AND timestamp = $4::timestamp);", l.Device, l.Latitude, l.Longitude, l.Timestamp); err != nil {
		return fmt.Errorf("connection: %w", err)
	}

	return nil
}

func (db *DB) CurrentLocations() ([]location.Location, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	rows, err := db.db.Query("SELECT DISTINCT ON (device) * FROM locations WHERE timestamp > (NOW() - INTERVAL '30 minutes') ORDER BY device, timestamp DESC;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []location.Location
	for rows.Next() {
		var loc location.Location
		if err := rows.Scan(&loc.Device, &loc.Latitude, &loc.Longitude, &loc.Timestamp); err != nil {
			return locations, err
		}
		locations = append(locations, loc)
	}

	if err = rows.Err(); err != nil {
		return locations, err
	}
	return locations, nil
}

func (db *DB) TimespanLocations(start, end time.Time) ([]location.Location, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	rows, err := db.db.Query("SELECT * FROM locations WHERE timestamp > $1 AND timestamp < $2 ORDER BY timestamp DESC;", start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []location.Location
	for rows.Next() {
		var loc location.Location
		if err := rows.Scan(&loc.Device, &loc.Latitude, &loc.Longitude, &loc.Timestamp); err != nil {
			return locations, err
		}
		locations = append(locations, loc)
	}

	if err = rows.Err(); err != nil {
		return locations, err
	}
	log.Printf("%d locations from %s till %s", len(locations), start.String(), end.String())
	return locations, nil
}

func (db *DB) DeviceLocations(device string) ([]location.Location, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	rows, err := db.db.Query("SELECT latitude, longitude, timestamp FROM locations WHERE device = $1 ORDER BY timestamp;", device)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []location.Location
	for rows.Next() {
		var loc location.Location
		if err := rows.Scan(&loc.Latitude, &loc.Longitude, &loc.Timestamp); err != nil {
			return locations, err
		}
		locations = append(locations, loc)
	}

	if err = rows.Err(); err != nil {
		return locations, err
	}
	return locations, nil
}

func (db *DB) checkDynamicDeviceID(device string) (bool, error) {
	var exists bool
	err := db.db.QueryRow("SELECT EXISTS (SELECT 1 FROM devices WHERE device = $1);", device).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (db *DB) Close() {
	db.db.Close()
}
