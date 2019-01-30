package data

import (
	"errors"
	"sync"
)

var (
	errValueNotFound = errors.New("Couldn't find")
)

// DatastoreMap is ... the datastore for this program (in memory).
// It implements `Datastore` and is thread-safe.
type DatastoreMap struct {
	m  map[string]*Data
	mu sync.Mutex
}

// Get looks for an element in the map.
func (ds *DatastoreMap) Get(websiteURL string, sessionID string) (*Data, bool, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	d, ok := ds.m[getStoreKey(websiteURL, sessionID)]
	return d, ok, nil
}

// Store adds/replaces the value on the specified key to the map.
func (ds *DatastoreMap) Store(websiteURL string, sessionID string, val *Data) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.m[getStoreKey(websiteURL, sessionID)] = val
	return nil
}

func getStoreKey(websiteURL string, sessionID string) string {
	return websiteURL + "/" + sessionID
}
