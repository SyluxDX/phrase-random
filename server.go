package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"./utils"
)

type phrase struct {
	Phrase string `json:"Message"`
}

// Configs server configurations
var (
	Configs  utils.Configurations
	Template *template.Template
)

func generatePhrase() string {
	// generate seed based on minutes
	seed, _ := strconv.ParseInt(time.Now().Format("020120061504"), 10, 64)
	rand.Seed(seed)
	// read file
	content, _ := ioutil.ReadFile(Configs.PhraseFile)
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	index := rand.Intn(len(lines))
	return lines[index]
}

func generateJSON() phrase {
	body := phrase{Phrase: generatePhrase()}
	return body
}

func checkAccept(data []string) bool {
	for _, item := range data {
		if strings.ToLower(item) == "application/json" {
			return true
		}
	}
	return false
}

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.Header["Accept"] != nil {
		if checkAccept(r.Header["Accept"]) {
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(generateJSON())
		} else {
			w.Header().Add("Content-Type", "text/html")
			Template.Execute(w, generateJSON())
		}
	} else {
		w.Header().Add("Content-Type", "text/html")
		Template.Execute(w, generateJSON())
	}
	log.Printf("Endpoint Hit from: %v\n", r.RemoteAddr)
}

func main() {
	Configs = utils.GetConfigs()
	Template, _ = template.ParseFiles(Configs.HTMLFile)

	log.Printf("Serving API on %[1]s port %[2]d (http://%[1]s:%[2]d/)\n", Configs.ServerURL, Configs.ServerPort)
	// Setup pages
	http.HandleFunc("/", homePage)
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
	url := fmt.Sprintf("%s:%d", Configs.ServerURL, Configs.ServerPort)
	log.Fatal(http.ListenAndServe(url, nil))
}
