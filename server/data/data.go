package data

// Data is the structure that holds the information about what the user is doing in the page.
// This will be "built up" over time, until the user presses the submit button.
type Data struct {
	WebsiteURL         string
	SessionID          string
	ResizeFrom         Dimension
	ResizeTo           Dimension
	CopyAndPaste       map[string]bool // map[fieldId]true
	FormCompletionTime int             // Seconds
}

// Dimension is the structure that holds the user page's dimensions (w x h).
type Dimension struct {
	Width  string `json:"width"`
	Height string `json:"height"`
}

// New receives a website URL and session ID (the only two required params for this)
// and returns the created data object ref, whilst adding it to the data store.
func New(websiteURL string, sessionID string) *Data {
	d := &Data{
		WebsiteURL: websiteURL,
		SessionID:  sessionID,
	}

	Ds.Store(d.WebsiteURL, d.SessionID, d)
	return d
}
