package main

import (
	"context"
	"fmt"
	//"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/t-hale/calendar/lib"
	"google.golang.org/api/calendar/v3"
	"log"
	"net/http"
	"os"
)

type handler struct{}

var (
	calendarService *calendar.Service
	httpHandler     handler
)

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/create":
		createCalendar(w, r)
	case "/delete":
		deleteCalendar(w, r)
	case "/list":
		listCalendars(w, r)
	case "/sync":
		syncCalendar(w, r)
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

//func entrypoint(w http.ResponseWriter, r *http.Request) {
//	cmd := &cli.Command{
//		Commands: []*cli.Command{
//			{
//				Name:    "create",
//				Aliases: []string{"a"},
//				Usage:   "create a shared calendar",
//				Action: func(ctx context.Context, cmd *cli.Command) error {
//					createCalendar(w, r)
//					return nil
//				},
//			},
//			{
//				Name:    "delete",
//				Aliases: []string{"c"},
//				Usage:   "delete a shared calendar",
//				Action: func(ctx context.Context, cmd *cli.Command) error {
//					deleteCalendar(w, r)
//					return nil
//				},
//			},
//			{
//				Name:    "list",
//				Aliases: []string{"a"},
//				Usage:   "list all shared calendars",
//				Action: func(ctx context.Context, cmd *cli.Command) error {
//					listCalendars(w, r)
//					return nil
//				},
//			},
//			{
//				Name:    "sync",
//				Aliases: []string{"c"},
//				Usage:   "sync a shared calendar",
//				Action: func(ctx context.Context, cmd *cli.Command) error {
//					syncCalendar(w, r)
//					return nil
//				},
//			},
//		},
//	}
//
//	if err := cmd.Run(context.Background(), os.Args); err != nil {
//		log.Fatal(err)
//	}
//}

func createCalendar(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var err error
	calendarService, err = calendar.NewService(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve calendar client: %v", err), http.StatusInternalServerError)
	}

	calendarId, err := lib.CreateSharedCalendar(calendarService, "not-the-primary-calendar")
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to create shared calendar: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("Created shared calendar: %s\n", calendarId)
}

func deleteCalendar(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var err error
	calendarService, err = calendar.NewService(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve calendar client: %v", err), http.StatusInternalServerError)
	}

	err = lib.DeleteCalendar(calendarService, "not-the-primary-calendar")
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to delete shared calendar: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("Deleted shared calendar: %s\n", "not-the-primary-calendar")
}

func listCalendars(w http.ResponseWriter, r *http.Request) {
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
	for _, cal := range calendars {
		log.Printf("Calendar: %+v\n", cal)
	}
}

func syncCalendar(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var err error
	calendarService, err = calendar.NewService(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve calendar client: %v", err), http.StatusInternalServerError)
	}
}
