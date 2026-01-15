package helloworld

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/thale/family-cal/lib"
	"google.golang.org/api/calendar/v3"
)

var (
	calendarService *calendar.Service
)

func init() {
	functions.HTTP("CreateSharedCalendar", createSharedCalendar)
}

func createSharedCalendar(w http.ResponseWriter, r *http.Request) {
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
