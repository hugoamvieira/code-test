package api

import "github.com/hugoamvieira/code-test/server/data"

// Every possible event will be listed here. This is done on purpose
// as I believe if an event is to be parsed, it must be explicitly declared here.
// Avoids strange undocumented behaviour.

type copyAndPasteEvent struct {
	WebsiteURL string `json:"websiteURL"`
	SessionID  string `json:"sessionID"`
	InputID    string `json:"inputID"`
}

type resizePageEvent struct {
	WebsiteURL string         `json:"websiteURL"`
	SessionID  string         `json:"sessionID"`
	ResizeFrom data.Dimension `json:"resizeFrom"`
	ResizeTo   data.Dimension `json:"resizeTo"`
}

type timeTakenEvent struct {
	WebsiteURL string `json:"websiteURL"`
	SessionID  string `json:"sessionID"`
	TimeTaken  int64  `json:"time"` // Seconds
}

type newSessionRequest struct {
	WebsiteURL string `json:"websiteURL"`
}

type newSessionResponse struct {
	SessionID string `json:"sessionID"`
}
