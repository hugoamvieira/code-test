package api

import (
	"net/url"

	"github.com/hugoamvieira/code-test/server/data"
)

// Every possible event will be listed here. This is done on purpose
// as I believe if an event is to be parsed, it must be explicitly declared here.
// Avoids strange undocumented behaviour.

type copyPasteEvent struct {
	WebsiteURL string `json:"websiteURL"`
	SessionID  string `json:"sessionID"`
	InputID    string `json:"inputID"`
}

func (cpe *copyPasteEvent) Valid() (bool, error) {
	// A session must already exist
	ok, err := validURLAndSession(cpe.WebsiteURL, cpe.SessionID)
	if err != nil {
		return false, err
	}
	return cpe.InputID != "" && ok, nil
}

type resizePageEvent struct {
	WebsiteURL string         `json:"websiteURL"`
	SessionID  string         `json:"sessionID"`
	ResizeFrom data.Dimension `json:"resizeFrom"`
	ResizeTo   data.Dimension `json:"resizeTo"`
}

func (rpe *resizePageEvent) Valid() (bool, error) {
	// A session must already exist
	ok, err := validURLAndSession(rpe.WebsiteURL, rpe.SessionID)
	if err != nil {
		return false, err
	}

	// Maybe validate if we can parse these into ints instead of this?
	validResizeFrom := rpe.ResizeFrom.Width != "" && rpe.ResizeFrom.Height != ""
	validResizeTo := rpe.ResizeTo.Width != "" && rpe.ResizeTo.Height != ""

	return validResizeTo && validResizeFrom && ok, nil
}

type timeTakenEvent struct {
	WebsiteURL string `json:"websiteURL"`
	SessionID  string `json:"sessionID"`
	TimeTaken  int    `json:"timeSeconds"` // Seconds
}

func (tte *timeTakenEvent) Valid() (bool, error) {
	// A session must already exist
	ok, err := validURLAndSession(tte.WebsiteURL, tte.SessionID)
	if err != nil {
		return false, err
	}

	return (tte.TimeTaken > 0) && ok, nil
}

type newSessionRequest struct {
	WebsiteURL string `json:"websiteURL"`
}

func (nsr *newSessionRequest) Valid() bool {
	// This is best-effort. Validating URLs is crazy difficult (too much ambiguity!)
	// Do we care that http://xyz.com and https://xyz.com are two separate websites in this system?
	// Do we care that http://xyz.com/ and http://xyz.com are also two separate websites?
	// More work would be required to make this better.
	// Thankfully, we have a single place where we can define an event's validity!
	// Yay for good design :P
	_, err := url.Parse(nsr.WebsiteURL)
	if err != nil {
		return false
	}
	return true
}

type newSessionResponse struct {
	SessionID string `json:"sessionID"`
}

func validURLAndSession(url string, session string) (bool, error) {
	_, ok, err := data.Ds.Get(url, session)
	return ok, err
}
