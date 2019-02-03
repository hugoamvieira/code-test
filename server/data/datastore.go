package data

var Ds Datastorer

func init() {
	Ds = &DatastoreMap{
		m: make(map[string]*Data),
	}
}

// Datastorer is the interface that defines the line between the application
// context and the datastore (currently, an in-memory Go map)
type Datastorer interface {
	Get(websiteURL string, sessionID string) (*Data, bool, error)
	Store(websiteURL string, sessionID string, val *Data) error
	Mutate(websiteURL string, sessionID string, newData *Data) (*Data, error)
}
