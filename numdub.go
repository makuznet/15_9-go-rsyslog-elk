package main

import (
	"fmt"
	"net/http"
	"strconv"

	// Phuslu log lib is in use. https://github.com/phuslu/log
	"github.com/phuslu/log"
)

func main() {

	// Greeting message for console users
	fmt.Println("This api doubles the entered integer number.\nEnter \"curl -X GET http://127.0.0.1:8080/v1/numdub/<your_number>\"")

	// Run numDub function when receiving ip:port/v1/numdub/ from http
	http.HandleFunc("/v1/numdub/", numDub)

	// the tcp port http server is listening to on all interfaces
	log.Fatal().Err(http.ListenAndServe(":8080", nil))
}

// Double the number and log
func numDub(w http.ResponseWriter, r *http.Request) {

	// Greeting message for web users
	fmt.Fprintf(w, "This api doubles the entered integer number.\nEnter \"curl -X GET http://127.0.0.1:8080/v1/numdub/<your_number>\"\n")

	// read url from http, cut url, leave entered string value, convert it to integer, assign to num var
	num, _ := strconv.Atoi(r.URL.Path[11:])

	if r.Method != "GET" {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}

	// Log enterd number
	log.Info().Msgf("number %d", num)

	// ptint out doubled number to http
	fmt.Fprintf(w, "doubled number: %d\n", num*2)

}
