package routes

import (
	"encoding/json"
	"net/http"
	"log"
	"sync"

	"qualiopi-backend/database"

	"github.com/gorilla/mux"
)

var (
	formatMu   sync.Mutex
	formatList []database.Format
)

func RegisterRoutes(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/formats", postFormat).Methods("POST")
	api.HandleFunc("/formats", getFormat).Methods("GET")
	api.HandleFunc("/formats/{id}/days", getDays).Methods("GET")
	api.HandleFunc("/formats/{id}/presents", getPresents).Methods("GET")
    api.HandleFunc("/days", addDay).Methods("POST")
    api.HandleFunc("/presents", addPresent).Methods("POST")
    api.HandleFunc("/presents/{id}", updatePresent).Methods("PATCH")
    api.HandleFunc("/formats/full", saveFull).Methods("POST")
    api.HandleFunc("/formats/{id}", getOneFormat).Methods("GET")
    api.HandleFunc("/formats/{id}", updateFormat).Methods("PUT")
    api.HandleFunc("/formats/{id}", deleteFormat).Methods("DELETE")

}

func postFormat(w http.ResponseWriter, r *http.Request) {
	var format database.Format

	if err := json.NewDecoder(r.Body).Decode(&format); err != nil {
		log.Println("decode format:", err)
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	formatMu.Lock()
	formatList = append(formatList, format)
	formatMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(format)
}

func getFormat(w http.ResponseWriter, r *http.Request) {
	formatMu.Lock()
	defer formatMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if formatList == nil {
		json.NewEncoder(w).Encode([]database.Format{})
		return
	}
	json.NewEncoder(w).Encode(formatList)
}

func getDays(w http.ResponseWriter, r *http.Request) {
	formatID := mux.Vars(r)["id"]
	database.Planner.Mu.RLock()
	defer database.Planner.Mu.RUnlock()

	out := []database.Day{}
	for _, d := range database.Planner.Days {
		if d.FormatID == formatID {
			out = append(out, d)
		}
	}
	json.NewEncoder(w).Encode(out)
}

func getPresents(w http.ResponseWriter, r *http.Request) {
	formatID := mux.Vars(r)["id"]
	database.Planner.Mu.RLock()
	defer database.Planner.Mu.RUnlock()

	out := []database.Present{}
	for _, p := range database.Planner.Presents {
		if d, ok := database.Planner.Days[p.DayID]; ok && d.FormatID == formatID {
			out = append(out, p)
		}
	}
	json.NewEncoder(w).Encode(out)
}

func addDay(w http.ResponseWriter, r *http.Request) {
	var d database.Day
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil || d.ID == "" {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	database.Planner.Mu.Lock()
	database.Planner.Days[d.ID] = d
	database.Planner.Mu.Unlock()
	w.WriteHeader(http.StatusCreated)
}

func addPresent(w http.ResponseWriter, r *http.Request) {
	var p database.Present
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil || p.ID == "" {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	database.Planner.Mu.Lock()
	database.Planner.Presents[p.ID] = p
	database.Planner.Mu.Unlock()
	w.WriteHeader(http.StatusCreated)
}

func updatePresent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var body struct {
		Present string `json:"present"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	database.Planner.Mu.Lock()
	defer database.Planner.Mu.Unlock()
	p, ok := database.Planner.Presents[id]
	if !ok {
		http.NotFound(w, r)
		return
	}
	p.Present = body.Present
	database.Planner.Presents[id] = p
	w.WriteHeader(http.StatusNoContent)
}

type fullFormat struct {
	Format   database.Format    `json:"format"`
	Days     []database.Day     `json:"days"`
	Presents []database.Present `json:"presents"`
}

func saveFull(w http.ResponseWriter, r *http.Request) {
	var body fullFormat
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("decode full:", err)
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	formatMu.Lock()
	formatList = append(formatList, body.Format)
	formatMu.Unlock()

	database.Planner.Mu.Lock()
	for _, d := range body.Days {
		database.Planner.Days[d.ID] = d
	}
	for _, p := range body.Presents {
		database.Planner.Presents[p.ID] = p
	}
	database.Planner.Mu.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func getOneFormat(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	formatMu.Lock()
	defer formatMu.Unlock()

	for _, f := range formatList {
		if f.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(f)
			return
		}
	}
	http.NotFound(w, r)
}

func updateFormat(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var f database.Format
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	f.ID = id // the URL wins over the body

	formatMu.Lock()
	defer formatMu.Unlock()

	for i := range formatList {
		if formatList[i].ID == id {
			formatList[i] = f
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}

func deleteFormat(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	formatMu.Lock()
	found := false
	for i := range formatList {
		if formatList[i].ID == id {
			formatList = append(formatList[:i], formatList[i+1:]...)
			found = true
			break
		}
	}
	formatMu.Unlock()

	if !found {
		http.NotFound(w, r)
		return
	}

	database.Planner.Mu.Lock()
	for dayID, d := range database.Planner.Days {
		if d.FormatID != id {
			continue
		}
		for pID, p := range database.Planner.Presents {
			if p.DayID == dayID {
				delete(database.Planner.Presents, pID)
			}
		}
		delete(database.Planner.Days, dayID)
	}
	database.Planner.Mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}