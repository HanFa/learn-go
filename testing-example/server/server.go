package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/double", doubleHandler)
	log.Fatalln(http.ListenAndServe(":4000", nil))
}

func doubleHandler(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("v")

	if text == "" {
		http.Error(w, "missing value", http.StatusBadRequest)
		return
	}

	v, err := strconv.Atoi(text)
	if err != nil {
		http.Error(w, "not a number: "+text, http.StatusBadRequest)
		return
	}

	if _, err = fmt.Fprintln(w, v*2); err != nil {
		http.Error(w, "cannot write to response", http.StatusBadRequest)
		return
	}

}
