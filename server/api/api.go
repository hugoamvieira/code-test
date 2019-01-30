package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hugoamvieira/code-test/server/data"
)

// API wraps Go's HTTP server. I've created it so it's physically and conceptually
// separated from the rest of the code and that so that any future API modifications
// are easier to do without changing other pieces of code (eg: Adding new things to the API struct)
type API struct {
	srv  *http.Server
	rand *rand.Rand
	sg   *sessionGen
}

const (
	errInvalidMethod   = `{"error": "Invalid method"}`
	errReadRequestBody = `{"error":"Couldn't read request body"}`
	errIncorrectJSON   = `{"error": "Malformed request body"}`
	errInvalidURL      = `{"error": "Invalid URL"}`
	errInternalServer  = `{"error": "Internal Server Error"}`
)

// New returns a new API object with a Go http server and a new serve mux with the
// API routes already defined.
func New(addr string) *API {
	a := &API{}

	m := http.NewServeMux()
	m.HandleFunc("/new_session", a.handleNewSession)
	m.HandleFunc("/new_resize_event", a.handleResizeEvent)
	m.HandleFunc("/new_cp_event", a.handleCopyPasteEvent)
	m.HandleFunc("/new_time_taken_event", a.handleTimeTakenEvent)

	a.srv = &http.Server{
		Addr:    addr,
		Handler: m,
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	a.rand = r

	a.sg = &sessionGen{
		sessions: make(map[string]bool),
	}

	return a
}

// Start starts the API, listening on all routes
func (a *API) Start() error {
	return a.srv.ListenAndServe()
}

func (a *API) handleNewSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errInvalidMethod, http.StatusMethodNotAllowed)
		return
	}
	fmt.Println(r.Host)

	rd := bufio.NewReader(r.Body)
	defer r.Body.Close()

	var bodyBytes []byte
	_, err := rd.Read(bodyBytes)
	if err != nil {
		http.Error(w, errReadRequestBody, http.StatusBadRequest)
		return
	}

	var ns newSessionRequest
	err = json.Unmarshal(bodyBytes, &ns)
	if err != nil {
		http.Error(w, errIncorrectJSON, http.StatusBadRequest)
		return
	}

	_, err = url.Parse(ns.WebsiteURL)
	if err != nil {
		log.Println("Failed to parse URL | Error:", err)
		http.Error(w, errInvalidURL, http.StatusBadRequest)
		return
	}

	// Create and store data object
	d := data.New(ns.WebsiteURL, a.generateSessionID())

	resp := newSessionResponse{
		SessionID: d.SessionID,
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Println("Failed to marshal response to JSON | Error:", err)
		http.Error(w, errInternalServer, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(respBytes)
	if err != nil {
		log.Println("Failed to write resp bytes to wire | Error:", err)
		http.Error(w, errInternalServer, http.StatusInternalServerError)
		return
	}
}

func (a *API) handleResizeEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errInvalidMethod, http.StatusMethodNotAllowed)
		return
	}
}

func (a *API) handleCopyPasteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errInvalidMethod, http.StatusMethodNotAllowed)
		return
	}
}

func (a *API) handleTimeTakenEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errInvalidMethod, http.StatusMethodNotAllowed)
		return
	}
}

func (a *API) generateSessionID() string {
	// Admittedly this is not great, however:
	// Since session IDs are also bounded by the website, I think we can get away with this
	// (in terms of uniqueness).
	//
	// Since I cannot use external libraries, the alternatives would be having an atomically incremental
	// int64 counter or create something like a k-sortable ID library.
	// Regarding the latter, it's just too time-consuming and frankly I'm not strong on crypto, so I opted to not go there.
	// Regarding the former, we wouldn't have any collisions but then we'd be permanently stuck
	// to 2^63-1 sessions unless we reset the counter at some arbitrary point, and that sounds not very pleasant.
	//
	// I chose to go with this approach because, even though it may generate more collisions,
	// we also have control of when the sessions end (when the user clicks the submit button),
	// so we can delete accordingly, which means we have _some_ control over collisions (it isn't just a
	// forever incrementing map). This means the biggest problem with the atomic counter is solved.
	// The downside to this approach is that, say there's a lot of slow users and a lot of traffic,
	// this loop could run for a long time until it finds a suitable sessionID.
	// Pratically, since this service won't have either, I feel like this is a good compromise.
	for {
		sessionID := strconv.FormatInt(a.rand.Int63(), 10)
		if ok := a.sg.Get(sessionID); !ok {
			a.sg.Set(sessionID)
			return sessionID
		}
	}
}
