package api

import "sync"

type sessionGen struct {
	sessions map[string]bool
	mu       sync.Mutex
}

func (sg *sessionGen) Get(sessionID string) bool {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	_, ok := sg.sessions[sessionID]
	return ok
}

func (sg *sessionGen) Set(sessionID string) {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	sg.sessions[sessionID] = true
}

func (sg *sessionGen) Delete(sessionID string) {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	delete(sg.sessions, sessionID)
}
