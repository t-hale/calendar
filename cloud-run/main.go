package cloud_run

import (
	"context"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"log"
	"net/http"
	"os"

	"github.com/thale/family-cal/lib"
	"github.com/urfave/cli/v3"
	"google.golang.org/api/calendar/v3"
)

var (
	calendarService *calendar.Service
)

func init() {
	functions.HTTP("calendar", entrypoint)
}

func entrypoint(w http.ResponseWriter, r *http.Request) {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:    "create",
				Aliases: []string{"a"},
				Usage:   "create a shared calendar",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					createCalendar(w, r)
					return nil
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"c"},
				Usage:   "delete a shared calendar",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					deleteCalendar(w, r)
					return nil
				},
			},
			{
				Name:    "list",
				Aliases: []string{"a"},
				Usage:   "list all shared calendars",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					listCalendars(w, r)
					return nil
				},
			},
			{
				Name:    "sync",
				Aliases: []string{"c"},
				Usage:   "sync a shared calendar",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					syncCalendar(w, r)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

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
