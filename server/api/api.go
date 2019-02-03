package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/hugoamvieira/code-test/server/data"
	"github.com/hugoamvieira/code-test/server/hash"
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
	errInvalidMethod      = `{"error": "Invalid method"}`
	errReadRequestBody    = `{"error":"Couldn't read request body"}`
	errIncorrectJSON      = `{"error": "Malformed request body"}`
	errInvalidRequest     = `{"error": "Invalid Request Body"}`
	errInternalServer     = `{"error": "Internal Server Error"}`
	errSessionNonExistent = `{"error": "Session doesn't exist"}`
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

func (a *API) writeOptions(w http.ResponseWriter) {
	a.setCorsHeaders(w)
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
}

func (a *API) setCorsHeaders(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*") // Don't do this in prod lol
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Add("Access-Control-Allow-Methods", "OPTIONS, POST")
}

func (a *API) handleNewSession(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		a.writeOptions(w)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, errInvalidMethod, http.StatusMethodNotAllowed)
		return
	}
	a.setCorsHeaders(w)
	w.Header().Add("Content-Type", "application/json")

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read request body | Error:", err)
		http.Error(w, errReadRequestBody, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var nsr newSessionRequest
	err = json.Unmarshal(bodyBytes, &nsr)
	if err != nil {
		log.Println("Failed to unmarshal JSON to struct | Error:", err)
		http.Error(w, errIncorrectJSON, http.StatusBadRequest)
		return
	}

	if !nsr.Valid() {
		http.Error(w, errInvalidRequest, http.StatusBadRequest)
		return
	}

	// Create and store data object
	sessionID := a.generateSessionID()
	if sessionID == "" {
		log.Printf("Couldn't generate a session ID for user in %v. Timed-out", nsr.WebsiteURL)
		http.Error(w, errInternalServer, http.StatusInternalServerError)
		return
	}

	d := data.New(nsr.WebsiteURL, sessionID)
	log.Printf("Data @ handleNewSession\n%#+v", d)
	log.Printf("Hash of %v: %v", d.WebsiteURL, hash.New(d.WebsiteURL))

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
	if r.Method == http.MethodOptions {
		a.writeOptions(w)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, errInvalidMethod, http.StatusMethodNotAllowed)
		return
	}
	a.setCorsHeaders(w)

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, errReadRequestBody, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var rpe resizePageEvent
	err = json.Unmarshal(bodyBytes, &rpe)
	if err != nil {
		log.Println("Failed to unmarshal JSON to struct | Error:", err)
		http.Error(w, errIncorrectJSON, http.StatusBadRequest)
		return
	}

	valid, err := rpe.Valid()
	if err != nil {
		log.Println("Failed to determine if request was valid (issues w/ datastore) | Error:", err)
		http.Error(w, errInternalServer, http.StatusInternalServerError)
		return
	}
	if !valid {
		http.Error(w, errInvalidRequest, http.StatusBadRequest)
		return
	}

	d := &data.Data{
		ResizeFrom: rpe.ResizeFrom,
		ResizeTo:   rpe.ResizeTo,
	}

	newData, err := data.Ds.Mutate(rpe.WebsiteURL, rpe.SessionID, d)
	if err != nil {
		log.Println("Error mutating data | Error:", err)
		http.Error(w, errInternalServer, http.StatusInternalServerError)
		return
	}

	log.Printf("Data @ handleResizeEvent\n%#+v", newData)
}

func (a *API) handleCopyPasteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		a.writeOptions(w)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, errInvalidMethod, http.StatusMethodNotAllowed)
		return
	}

	a.setCorsHeaders(w)

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, errReadRequestBody, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var cpe copyPasteEvent
	err = json.Unmarshal(bodyBytes, &cpe)
	if err != nil {
		log.Println("Failed to unmarshal JSON to struct | Error:", err)
		http.Error(w, errIncorrectJSON, http.StatusBadRequest)
		return
	}

	valid, err := cpe.Valid()
	if err != nil {
		log.Println("Failed to determine if request was valid (issues w/ datastore) | Error:", err)
		http.Error(w, errInternalServer, http.StatusInternalServerError)
		return
	}
	if !valid {
		http.Error(w, errInvalidRequest, http.StatusBadRequest)
		return
	}

	d := &data.Data{
		CopyAndPaste: map[string]bool{
			cpe.InputID: true,
		},
	}

	newData, err := data.Ds.Mutate(cpe.WebsiteURL, cpe.SessionID, d)
	if err != nil {
		log.Println("Error mutating data | Error:", err)
		http.Error(w, errInternalServer, http.StatusInternalServerError)
		return
	}

	log.Printf("Data @ handleCopyPasteEvent\n%#+v", newData)
}

func (a *API) handleTimeTakenEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		a.writeOptions(w)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, errInvalidMethod, http.StatusMethodNotAllowed)
		return
	}
	a.setCorsHeaders(w)

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, errReadRequestBody, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var tte timeTakenEvent
	err = json.Unmarshal(bodyBytes, &tte)
	if err != nil {
		log.Println("Failed to unmarshal JSON to struct | Error:", err)
		http.Error(w, errIncorrectJSON, http.StatusBadRequest)
		return
	}

	valid, err := tte.Valid()
	if err != nil {
		log.Println("Failed to determine if request was valid (issues w/ datastore) | Error:", err)
		http.Error(w, errInternalServer, http.StatusInternalServerError)
		return
	}
	if !valid {
		http.Error(w, errInvalidRequest, http.StatusBadRequest)
		return
	}

	d := &data.Data{
		FormCompletionTime: tte.TimeTaken,
	}

	newData, err := data.Ds.Mutate(tte.WebsiteURL, tte.SessionID, d)
	if err != nil {
		log.Println("Error mutating data | Error:", err)
		http.Error(w, errInternalServer, http.StatusInternalServerError)
		return
	}

	log.Printf("Data @ handleTimeTakenEvent\n%#+v", newData)
}

func (a *API) generateSessionID() string {
	// Admittedly this is not great, however:
	// Since session IDs are also bounded by the website, I think we can get away with this
	// (in terms of uniqueness).
	//
	// Since I cannot use external libraries, the options I thought of are:
	// 1. Having an atomically incremented int64 counter;
	// 2. Create something like a K-Sortable ID library;
	// 3. Keep a store of all the used sessions and just use rand.Int63.
	//
	// Regarding 1., we wouldn't have any collisions but then we'd be permanently stuck
	// to 2^63-1 sessions unless we reset the counter at some arbitrary point, and that sounds not very pleasant.
	//
	// Regarding 2., frankly I'm not too strong on crypto so it'd just too time-consuming.
	//
	// I chose to go with 3. because, even though it may generate collisions, it only really starts becoming
	// an issue when we've generated over 3.3 billion sessions or so, which is when the probability raises
	// to 25% (assuming Go's rand is not biased and each number has the same probability of being picked)
	// It also has the same problem as 1. does, however since we have control over how we calculate the actual value,
	// we can add more uniqueness over time. We also know when a session is finished so we could theoretically free it to be used
	// again, but our "primary keys" are website url and session ID, so if we did that there'd be a chance that data
	// could be replaced _if_ a user goes into the same website _and_ gets the same session ID, which is highly unlikely but
	// I didn't want to take that risk. We could add more uniqueness by, say, integrating the user's IP into the mix.
	//
	// The downside to this approach is that, say there's a lot of traffic (again, we're talking billions of accrued users),
	// this loop could start running for a long time until it finds a suitable sessionID.
	// Mathematically speaking, there's also a point where this loop would run forever (when the probability of collision is ~75% or so, so I've added
	// a time-out to cover that. Also, who wants to wait that long for a sessionID?
	for start := time.Now(); time.Since(start) < 5*time.Second; {
		sessionID := strconv.FormatInt(a.rand.Int63(), 10)
		if ok := a.sg.Get(sessionID); !ok {
			a.sg.Set(sessionID)
			return sessionID
		}
	}
	return ""
}
