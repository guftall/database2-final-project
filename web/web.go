package web

import (
	"fmt"
	"io/ioutil"
	"olympic/db"

	"log"
	"net/http"
	"encoding/json"
)

func Start(listeningUrl string) {
	log.Printf("Starting server(%s)\n", listeningUrl)
	http.HandleFunc("/olympics", olympicsHandler)
	http.HandleFunc("/athletes", athletesHandler)
	http.HandleFunc("/athlete", athletesHandlerDelete)
	log.Fatal(http.ListenAndServe(listeningUrl, nil))
}


func olympicsHandler(w http.ResponseWriter, r *http.Request) {
	olympics := db.ReadOlympics()
	data := encode(olympics)


	w.Write(data)
}

func athletesHandler(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	year := r.URL.Query().Get("year")
	country := r.URL.Query().Get("country")

	athletes := db.ReadAthletes(name, year, country)
	data := encode(athletes)

	a := encode(athletes[0])
	// println(athletes[0].Id)

	// println(len(athletes[0].Id))

	// println(data)

	var ath db.Athlete
	err := json.Unmarshal(a, &ath)
	_ = err

	// println(ath.Id)

	w.Write(data)
}

func athletesHandlerDelete(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		id := r.URL.Query().Get("id")

		
		athlete := db.ReadAthlete(id)
		data := encode(athlete)
		w.Write(data)
	} else if r.Method == "DELETE" {

	}
	return
	// decoder := json.NewDecoder(r.Body)


	// type bod struct {
	// 	Id []byte	`json:"id"`
	// }
	// var receivedBody bod
	// err := decoder.Decode(&receivedBody)
    // if err != nil {
    //     panic(err)
	// }
	// log.Println(receivedBody)
	// return
	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	fmt.Printf("body: %v", body)
	println(len(body))

	db.DeleteAthletes(body)

	w.WriteHeader(200)
}

func encode(i interface{}) []byte {

	jsonData, err := json.Marshal(i)
	if err != nil {
		log.Printf("unable to parse json: %s\n", err)
	}

	return jsonData
}