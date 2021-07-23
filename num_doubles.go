package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {

	fmt.Println("This api doubles the number. Enter curl -X GET http://127.0.0.1:8080/v1/numdoubles/<your_number>")

	http.HandleFunc("/v1/numdoubles/", numDoubles)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Double the number and log
func numDoubles(w http.ResponseWriter, r *http.Request) {

	num, _ := strconv.Atoi(r.URL.Path[15:])

	if r.Method != "GET" {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}
	fmt.Println(num)
	fmt.Fprintf(w, "Doubled number: %d\n", num*2)
}
