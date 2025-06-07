package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
)

var (
	mutex = &sync.Mutex{}
	file  = "runs.json"
)

func main() {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		if err := os.WriteFile(file, []byte("{}"), 0644); err != nil {
			panic(err)
		}
	}
	http.HandleFunc("/runs", handleRuns)
	panic(http.ListenAndServe(":7867", nil))
}

func handleRuns(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name parameter missing", http.StatusBadRequest)
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	data := make(map[string]int)
	if content, err := os.ReadFile(file); err == nil {
		json.Unmarshal(content, &data) // Ignore unmarshal errors for empty/invalid JSON
	}
	data[name]++
	if content, err := json.MarshalIndent(data, "", "\t"); err == nil {
		os.WriteFile(file, content, 0644)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"runs": data[name]})
}
