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
	Width  string
	Height string
}
