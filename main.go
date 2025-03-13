package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

var mutex = &sync.Mutex{}

func main() {
	if _, err := os.Stat("runs.json"); os.IsNotExist(err) {
		initialData := make(map[string]int)
		initialDataBytes, _ := json.Marshal(initialData)
		err := ioutil.WriteFile("runs.json", initialDataBytes, 0644)
		if err != nil {
			fmt.Println("Error creating data.json:", err)
			return
		}
	}
	http.HandleFunc("/runs", handleRuns)
	port := "7867"
	fmt.Println("Server listening on port", port)
	http.ListenAndServe(":"+port, nil)
}

func handleRuns(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name parameter missing", http.StatusBadRequest)
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	file, err := ioutil.ReadFile("runs.json")
	if err != nil && !os.IsNotExist(err) {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	data := make(map[string]int)
	if err := json.Unmarshal(file, &data); err != nil && !os.IsNotExist(err) {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if val, ok := data[name]; ok {
		data[name] = val + 1
	} else {
		data[name] = 1
	}
	file, err = json.MarshalIndent(data, "", "\t")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = ioutil.WriteFile("runs.json", file, 0644)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	response := map[string]int{"runs": data[name]}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
