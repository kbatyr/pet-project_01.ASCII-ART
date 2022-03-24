package main

import (
	"fmt"
	"log"
	"net/http"

	"web/server"
)

// process the main page of the site
// start a local server on port 8080
func main() {
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("./style/"))))
	http.HandleFunc("/", server.HomePageHandler)
	http.HandleFunc("/ascii-art", server.AsciiHandler)
	fmt.Println("Server is listening at :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
