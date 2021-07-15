package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Input struct {
	Username string `json:"username"`
	RepoName string `json:"reponame"`
}

func getDetails(w http.ResponseWriter, r *http.Request) {
	/*reading the corresponding input parameters
	1. githubusername.
	2. github reponame
	3. AccessToken available as env variable.
	*/
	log.Println("Hitting Api Endpoint for HB Operations")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var input Input
	_ = json.NewDecoder(r.Body).Decode(&input)
	username := input.Username
	reponame := input.RepoName
	log.Println("Username", username)
	log.Println("RepoName", reponame)
	log.Println("params", params)
	log.Printf("Getting Info for %s user for %s repository.\n", username, reponame)
	finalResponse := callFlow(reponame)
	json.NewEncoder(w).Encode(finalResponse)

}

func main() {
	// initiating router
	router := mux.NewRouter()
	router.HandleFunc("/getUser/", getDetails).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
