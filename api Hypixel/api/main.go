package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type AllStruct struct {
	Success bool   `json:"success"`
	exemple string `json:"lastUpdated"`
	Product map[string]struct {
		Info struct {
			attack     float64 `json:"attack"`
			defense    float64 `json:"defense"`
			magic      float64 `json:"magic"`
			difficulty float64 `json:"difficulty"`
		} `json:"info"`
	} `json:"products"`
}

type joueur struct {
	Success bool   `json:"success"`
	Cause   string `json:"cause"`
	Player  struct {
		ID          string `json:"_id"`
		UUID        string `json:"uuid"`
		DisplayName string `json:"displayname"`
	} `json:"player"`
}
type PlayerResponse struct {
	Player  map[string]interface{} `json:"player"`
	Cause   string                 `json:"cause"`
	Success bool                   `json:"success"`
}
type KeyInfoResponse struct {
	Record  map[string]interface{} `json:"record"`
	Cause   string                 `json:"cause"`
	Success bool                   `json:"success"`
}
type GuildIDResponse struct {
	Guild   string `json:"guild"`
	Cause   string `json:"cause"`
	Success bool   `json:"success"`
}
type GuildResponse struct {
	Guild   map[string]interface{} `json:"guild"`
	Cause   string                 `json:"cause"`
	Success bool                   `json:"success"`
}
type FriendsResponse struct {
	Records []map[string]interface{} `json:"records"`
	Cause   string                   `json:"cause"`
	Success bool                     `json:"success"`
}
type SessionResponse struct {
	Session map[string]interface{} `json:"session"`
	Cause   string                 `json:"cause"`
	Success bool                   `json:"success"`
}

var words string = ""

const (
	key  = "386e7185-cf0c-4a8b-8d13-6beb3c36fd24"
	Host = "localhost"
	Port = "4444"
)

var BaseURL string = "https://api.hypixel.net/player?key=" + key + "&name=neder_2503"

func main() {

	httpClient := http.Client{
		Timeout: time.Second * 8, // define timeout
	}

	//create template file

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatalln(err)
	}

	//create request
	req, err := http.NewRequest(http.MethodGet, BaseURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	//make api call
	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	//parse response
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	response := joueur{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//imageServer := http.FileServer(http.Dir("images"))
		//http.Handle("/images/", http.StripPrefix("/images/", imageServer))
		words = r.FormValue("w")
		tmpl = template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, response)
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(Host+":"+Port, nil)
}
