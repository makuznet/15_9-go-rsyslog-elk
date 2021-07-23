package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	// Logrus log lib is in use. https://github.com/sirupsen/logrus#readme
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {

	// Greeting message for console users
	fmt.Println("This api doubles the entered integer number. Enter curl -X GET http://127.0.0.1:8080/v1/numdoubles/<your_number>")

	// Creating a log file
	file, err := os.OpenFile("numdub.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// Keep the log file open while the app works
	defer file.Close()

	// Write log to a log file
	log.SetOutput(file)

	// The app adds this line every time it starts
	log.Println("The app has been launched")

	// Run numDoubles function when receiving ip:port/v1/numdoubles/ from http
	http.HandleFunc("/v1/numdoubles/", numDoubles)

	// the tcp port http server is listening to on all interfaces
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Double the number and log
func numDoubles(w http.ResponseWriter, r *http.Request) {

	// Greeting message for web users
	fmt.Fprintf(w, "This api doubles the entered integer number. http://127.0.0.1:8080/v1/numdoubles/<your_number>\n")

	// read url from http, cut url, leave entered string value, convert it to integer, assign to num var
	num, _ := strconv.Atoi(r.URL.Path[15:])

	if r.Method != "GET" {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}
	fmt.Println(num)

	// ptint out entered number to http
	fmt.Fprintf(w, "Doubled number: %d\n", num*2)

	// write entered number to numdub.log
	log.WithFields(log.Fields{
		"number": num,
	}).Info("Number entered")
}
