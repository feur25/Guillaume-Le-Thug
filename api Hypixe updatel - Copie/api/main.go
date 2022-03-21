package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Player struct {
	ID          string `json:"_id"`
	UUID        string `json:"uuid"`
	DisplayName string `json:"displayname"`
	Rank        string `json:"rank"`
	LastLogin   int64  `json:"lastLogin"`
	LastLogout  int64  `json:"lastLogout"`
	PlayerName  string `json:"playername"`
}

type Records struct {
	ID           string `json:"_id"`
	UUIDSender   string `json:"uuidSender"`
	UUIDReceiver string `json:"uuidReceiver"`
	Started      int64  `json:"started"`
}
type AllStruct struct {
	Success bool      `json:"success"`
	Cause   string    `json:"cause"`
	Player  Player    `json:"player"`
	Records []Records `json:"records"`
}

const (
	key  = "386e7185-cf0c-4a8b-8d13-6beb3c36fd24"
	Host = "localhost"
	Port = "5555"
)

var tmpl = template.Must(template.ParseFiles("index.html"))
var fs = http.FileServer(http.Dir("css"))

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
}

var words = ""
var power = ""
var BaseURL string = ""
var friend_URL string = ""
var player_URL string = ""
var uuid_friend = ""
var uuid_player = ""
var n int = 0
var text2 = ""

func find_player(w http.ResponseWriter, r *http.Request) {
	words = r.PostFormValue("name_player")

	httpClient := http.Client{
		Timeout: time.Second * 12,
	}
	response := AllStruct{}
	BaseURL = "https://api.hypixel.net/player?key=" + key + "&name=" + words
	//BaseURL = "https://api.hypixel.net/player?key=" + key + "&name=" + words

	/*var findfriendurl string = "https://api.hypixel.net/player?key=" + key + "&uuid=" + words*/

	req, err := http.NewRequest(http.MethodGet, BaseURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	if response.Player.LastLogin > response.Player.LastLogout {
		response.Player.PlayerName = "Online"
	}
	if response.Player.LastLogin < response.Player.LastLogout {
		response.Player.PlayerName = "Offline"
	}
	//BaseURL = "https://api.hypixel.net/friends?key=" + key + "&uuid=" + uuid_friend
	if len(response.Player.UUID) > 0 {
		uuid_friend = response.Player.UUID
	}
	if len(uuid_friend) > 0 {
		friend_URL = "https://api.hypixel.net/friends?key=" + key + "&uuid=" + uuid_friend
		log.Print(uuid_friend)
		rep, err := http.NewRequest(http.MethodGet, friend_URL, nil)
		if err != nil {
			log.Fatal(err)
		}
		rest, getErr := httpClient.Do(rep)
		if getErr != nil {
			log.Fatal(getErr)
		}
		if res.Body != nil {
			defer res.Body.Close()
		}

		bodys, readErr := ioutil.ReadAll(rest.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		jsonErrs := json.Unmarshal(bodys, &response)
		if jsonErrs != nil {
			log.Fatal(jsonErrs)
		}
	}
	/*next*/
	power = r.FormValue("other")
	player_URL = "https://api.hypixel.net/player?key=" + key + "&uuid=" + power
	rep, err := http.NewRequest(http.MethodGet, player_URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	rest, getErr := httpClient.Do(rep)
	if getErr != nil {
		log.Fatal(getErr)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	bodys, readErr := ioutil.ReadAll(rest.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErrs := json.Unmarshal(bodys, &response)
	if jsonErrs != nil {
		log.Fatal(jsonErrs)
	}

	tmpl.Execute(w, response)
	log.Print(response.Player.PlayerName)
}

func main() {
	/*http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		find_player(w, r)
	})*/
	/*http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		id := strings.ReplaceAll(r.URL.Path, "/", "")
		url := "https://api.hypixel.net/player?key=" + key + "&name=" + id
	})*/
	http.HandleFunc("/", find_player)
	print("Lancement de la page instancier sur : " + Host + ":" + Port)
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	/*http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))*/

	http.ListenAndServe(Host+":"+Port, nil)

}
