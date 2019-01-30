package data

var Ds Datastore

func init() {
	Ds = &DatastoreMap{
		m: make(map[string]*Data),
	}
}

// Datastore is the interface that defines the line between the application
// context and the datastore (currently, an in-memory Go map)
type Datastore interface {
	Get(websiteURL string, sessionID string) (*Data, bool, error)
	Store(websiteURL string, sessionID string, val *Data) error
}
