package data

import (
	"errors"
	"sync"
)

var (
	errValueNotFound = errors.New("Couldn't find value for key")
	errNilValue      = errors.New("Nil value has been passed")
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

// Mutate retrieves the old data and replaces it with a new piece of data if the old data doesn't
// contain a particular parameter (so, sort of a diff)
// Calling Mutate on a url/session ID combo that doesn't exist will end up in an error.
// At the end, it'll return the "diff-ed" object.
func (ds *DatastoreMap) Mutate(websiteURL string, sessionID string, newData *Data) (*Data, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if newData == nil {
		return nil, errNilValue
	}

	oldData, ok := ds.m[getStoreKey(websiteURL, sessionID)]
	if !ok {
		return nil, errValueNotFound
	}

	// Since the Go structs don't use pointers, we have to check zero values for everything... *sadface*
	oldDataHasResizeFrom := oldData.ResizeFrom.Height != "" && oldData.ResizeFrom.Width != ""
	newDataHasResizeFrom := newData.ResizeFrom.Height != "" && newData.ResizeFrom.Width != ""

	oldDataHasResizeTo := oldData.ResizeTo.Height != "" && oldData.ResizeTo.Width != ""
	newDataHasResizeTo := newData.ResizeTo.Height != "" && newData.ResizeTo.Width != ""

	if (!oldDataHasResizeFrom && newDataHasResizeFrom) && (!oldDataHasResizeTo && newDataHasResizeTo) {
		// These only make sense if they're replaced together, I think
		oldData.ResizeFrom = newData.ResizeFrom
		oldData.ResizeTo = newData.ResizeTo
	}

	if oldData.FormCompletionTime == 0 && newData.FormCompletionTime > 0 {
		oldData.FormCompletionTime = newData.FormCompletionTime
	}

	// Add & Replace data from new copy and paste map to the old one.
	// Replacing regardless of existence in the old map since any already existent value
	// will only be replaced by itself.
	// We also don't have to worry about this map growing large (and thus making this loop slower in worst case)
	// as it (almost) directly correlates to fields that people have to put things in per website!
	for k, v := range newData.CopyAndPaste {
		if _, ok := oldData.CopyAndPaste[k]; !ok {
			oldData.CopyAndPaste[k] = v
		}
	}

	return oldData, nil
}

func getStoreKey(websiteURL string, sessionID string) string {
	return websiteURL + "/" + sessionID
}
