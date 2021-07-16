package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	log.Printf("... starting ...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("server /")
		fmt.Fprintf(w, "Hello, %q \n", html.EscapeString(r.URL.Path))

		for _, e := range os.Environ() {
			pair := strings.SplitN(e, "=", 2)
			fmt.Fprintf(w, "%s : %s \n", pair[0], pair[1])
		}

		//fmt.Fprintf(w, "PROPERTY_0: %s \n", os.Getenv("PROPERTY_0"))
		//fmt.Fprintf(w, "PROPERTY_1: %s \n", os.Getenv("PROPERTY_1"))
		//fmt.Fprintf(w, "PROPERTY_2: %s \n", os.Getenv("PROPERTY_2"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
