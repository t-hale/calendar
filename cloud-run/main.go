package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	calendarpb "github.com/t-hale/calendar/gen"
	"github.com/t-hale/calendar/lib"
	"google.golang.org/api/calendar/v3"
)

type handler struct{}

var (
	calendarService *calendar.Service
	httpHandler     handler
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>%s</title>
</head>
<body>
    <h1>User Struct (Pretty Printed JSON)</h1>
    <pre><code>%s</code></pre>
</body>
</html>`

// ParseJSON is a generic function to decode the request body into type T.
func ParseJSON[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func prettyPrintHTML[T any](title string, data T, w http.ResponseWriter) {
	// Marshal the struct to a pretty-printed JSON byte slice.
	// The prefix "" and indent "\t" create the indentation.
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the HTML response, embedding the JSON string safely.
	fmt.Fprintf(w, htmlTemplate, title, template.HTMLEscapeString(string(jsonData)))
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/create":
		req, err := ParseJSON[calendarpb.CreateCalendarRequest](r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to parse request body: %v", err), http.StatusBadRequest)
			return
		}
		createCalendar(w, &req)
	case "/delete":
		req, err := ParseJSON[calendarpb.DeleteCalendarRequest](r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to parse request body: %v", err), http.StatusBadRequest)
			return
		}
		deleteCalendar(w, &req)
	case "/list":
		//TODO - Add parsing back once the ListCalendarRequest body isn't nil
		//req, err := ParseJSON[calendarpb.ListCalendarRequest](r)
		//if err != nil {
		//	http.Error(w, fmt.Sprintf("Unable to parse request body: %v", err), http.StatusBadRequest)
		//	return
		//}
		listCalendars(w, &calendarpb.ListCalendarRequest{})
	case "/sync":
		req, err := ParseJSON[calendarpb.SyncCalendarRequest](r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to parse request body: %v", err), http.StatusBadRequest)
			return
		}
		syncCalendar(w, &req)
	default:
		// Default action or 404
		http.Error(w, fmt.Sprintf("URL %s unsupported", r.URL.Path), http.StatusNotFound)
	}
}

func init() {
	httpHandler = handler{}
}

func main() {
	log.Printf("starting calendar server")

	// Determine port for HTTP service from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 if not set locally
		log.Printf("defaulting to port %s", port)
	}

	// Start the HTTP server
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, httpHandler); err != nil {
		log.Fatal(err)
	}
}

func createCalendar(w http.ResponseWriter, req *calendarpb.CreateCalendarRequest) {
	ctx := context.Background()
	var err error
	calendarService, err = calendar.NewService(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve calendar client: %v", err), http.StatusInternalServerError)
	}

	calendarId, err := lib.CreateSharedCalendar(calendarService, req.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to create shared calendar: %v", err), http.StatusInternalServerError)
		return
	}
	message := fmt.Sprintf("Created shared calendar %s with id : %s", req.Name, calendarId)
	log.Printf(message)
	w.Write([]byte(message))
}

func deleteCalendar(w http.ResponseWriter, req *calendarpb.DeleteCalendarRequest) {
	ctx := context.Background()
	var err error
	calendarService, err = calendar.NewService(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve calendar client: %v", err), http.StatusInternalServerError)
	}

	err = lib.DeleteCalendar(calendarService, req.CalendarId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to delete shared calendar: %v", err), http.StatusInternalServerError)
		return
	}
	message := fmt.Sprintf("Deleted shared calendar: %s\n", req.CalendarId)
	log.Printf(message)
	w.Write([]byte(message))
}

func listCalendars(w http.ResponseWriter, req *calendarpb.ListCalendarRequest) {
	ctx := context.Background()
	var err error
	calendarService, err = calendar.NewService(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve calendar client: %v", err), http.StatusInternalServerError)
	}
	calendars, err := lib.GetCalendars(calendarService)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve calendars: %v", err), http.StatusInternalServerError)
	}
	log.Printf("%+v", calendars)
	prettyPrintHTML("ListCalendarResponse", calendars, w)
}

func syncCalendar(w http.ResponseWriter, req *calendarpb.SyncCalendarRequest) {
	ctx := context.Background()
	var err error
	calendarService, err = calendar.NewService(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve calendar client: %v", err), http.StatusInternalServerError)
	}

	message := fmt.Sprintf("Successfully synced calendar %s", req.CalendarId)
	log.Printf(message)
	w.Write([]byte(message))
}
