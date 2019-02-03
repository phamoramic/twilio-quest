package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

// Response is the content for the smsHandler response
type Response struct {
	Message string
}

func smsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		if fromCountry, ok := r.PostForm["FromCountry"]; ok {
			resp := Response{
				Message: fmt.Sprintf("Hi! It looks like your phone number was born in %s", fromCountry[0]),
			}
			x, err := xml.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Cannot marshal response: %s", err)))
			}
			w.Header().Set("Content-Type", "application/xml")
			w.Write(x)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Request body does not contain FromCountry field."))
		}
		return
	case http.MethodGet:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("GET not implemented. Use POST instead."))
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Use POST method."))
		return
	}
}

func main() {
	http.HandleFunc("/sms", smsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
