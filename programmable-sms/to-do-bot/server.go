package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var items []string

// Response is the content for the smsHandler response
type Response struct {
	Message string
}

func smsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		if body, ok := r.PostForm["Body"]; ok {
			bodies := strings.Fields(body[0])
			command := strings.ToUpper(bodies[0])
			var resp Response
			switch command {
			case "ADD":
				items = append(items, strings.Join(bodies[1:], " "))
				resp = Response{"Item added."}
			case "LIST":
				var buf bytes.Buffer
				for i, item := range items {
					buf.WriteString(fmt.Sprintf("%d. %s\n", i+1, item))
				}
				resp = Response{
					Message: buf.String(),
				}
			case "REMOVE":
				itemNum, err := strconv.Atoi(bodies[1])
				var respStr string
				if err != nil {
					respStr = "Item to remove must be a #."
				} else {
					itemNum = itemNum - 1 // Slice indices start at 0
					items = append(items[:itemNum], items[itemNum+1:]...)
					respStr = "Item removed."
				}
				resp = Response{respStr}
			default:
				resp = Response{"No action performed. Valid commands: ADD <task>, LIST, REMOVE <#>."}
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
			w.Write([]byte("Request body does not contain Body field."))
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
