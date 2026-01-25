package lib

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/apognu/gocal"
	"google.golang.org/api/calendar/v3"
)

type CalendarEntry struct {
	Name string
	Url  string
}

func cleanupTempCalendars(calendarService *calendar.Service) error {
	allCalendars, err := GetCalendars(calendarService)
	if err != nil {
		return fmt.Errorf("unable to retrieve calendar list: %v", err)
	}

	for _, c := range allCalendars {
		if c.Summary == "temp-calendar" {
			log.Printf("Deleting temporary calendar: %s\n", c.Id)
			if err := calendarService.Calendars.Delete(c.Id).Do(); err != nil {
				return fmt.Errorf("error deleting temporary calendar: %v", err)
			}
		}
	}
	return nil
}

func clearCalendar(calendarService *calendar.Service, calendarId string) error {
	log.Printf("Clearing calendar %s\n", calendarId)
	events, err := calendarService.Events.List(calendarId).Do()
	if err != nil {
		return fmt.Errorf("error retrieving events from calendar: %v", err)
	}

	for _, event := range events.Items {
		log.Printf("Deleting event %s from calendar %s\n", event.Id, calendarId)
		if err := calendarService.Events.Delete(calendarId, event.Id).Do(); err != nil {
			return fmt.Errorf("error deleting event from calendar: %v", err)
		}
	}

	return nil
}

// CreateSharedCalendar creates or retrieves a shared calendar with the given name using the provided calendar service.
// If a matching calendar exists, it returns its ID. Otherwise, it creates a new calendar and returns its ID.
// Returns an error if the calendar cannot be retrieved or created.
func CreateSharedCalendar(calendarService *calendar.Service, calendarName string) (string, error) {

	allCalendars, err := GetCalendars(calendarService)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve calendar list: %v", err)
	}

	var sharedCalendar string
	calendarExists := false
	for _, c := range allCalendars {
		if c.SummaryOverride == "" {
			log.Printf("Found calendar: %+v\n", c.Summary)
		} else {
			log.Printf("Found calendar: %+v\n", c.SummaryOverride)
		}
		if c.Summary == calendarName || c.SummaryOverride == calendarName {
			log.Printf("Found existing shared calendar: %s(%s), skipping creation\n", calendarName, c.Id)
			sharedCalendar = c.Id
			calendarExists = true
			break
		}
	}

	if !calendarExists {
		sc, err := calendarService.Calendars.Insert(&calendar.Calendar{Summary: calendarName}).Do()

		if err != nil {
			return "", fmt.Errorf("unable to create shared calendar: %v", err)
		}

		log.Printf("Created shared calendar: %v\n", sc.Id)
		sharedCalendar = sc.Id
	}

	calendarListEntry := &calendar.CalendarListEntry{
		Summary:    calendarName,
		AccessRole: "owner",
		Id:         sharedCalendar,
	}
	cle, err := calendarService.CalendarList.Insert(calendarListEntry).Do()
	if err != nil {
		return "", fmt.Errorf("unable to add shared calendar to list of calendars: %v", err)
	}

	log.Printf("Updated shared calendar and added to calendar lists: %v\n", cle.Id)
	return sharedCalendar, nil
}

func DeleteCalendar(calendarService *calendar.Service, calendarId string) error {
	log.Printf("Deleting calendar %s\n", calendarId)
	if err := calendarService.Calendars.Delete(calendarId).Do(); err != nil {
		return fmt.Errorf("error deleting calendar: %v", err)
	}
	return nil
}

// GetCalendars retrieves all calendars available in the user's Google Calendar account using the provided service.
// Returns a list of calendar entries or an error if the operation fails.
func GetCalendars(calendarService *calendar.Service) ([]*calendar.CalendarListEntry, error) {
	calendars, err := calendarService.CalendarList.List().Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve calendar list: %v", err)
	}
	return calendars.Items, nil
}

// listCalendarEvents retrieves and logs the upcoming events from the specified calendar using the provided service.
// Takes a calendar service instance and a calendar ID as inputs. Returns an error if the events cannot be retrieved.
func listCalendarEvents(calendarService *calendar.Service, calendarId string) error {
	events, err := calendarService.Events.List(calendarId).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve next ten events: %v", err)
	}
	log.Println("Upcoming events:")
	for _, item := range events.Items {
		log.Printf("\t%v (%v)\n", item.Summary, item.Start.DateTime)
	}

	return nil
}

// ImportEventsToSharedCalendar imports a list of events into a specified shared calendar using the provided calendar service.
// It requires a calendar service instance, the target calendar ID where events will be imported, and a list of CalendarEntry objects of events.
// Returns an error if the import fails, or nil if all events are successfully imported.
func ImportEventsToSharedCalendar(calendarService *calendar.Service, calendarId string, calendarEntries []*CalendarEntry) error {
	for _, ce := range calendarEntries {
		// logic to import the calendar URL into the named calendar would go here
		log.Printf("Importing calendar %s: %s\n", ce.Name, ce.Url)

		now := time.Now()
		nextYear := now.AddDate(1, 0, 0)
		events, err := getICSCalendarEvents(ce.Url, now, nextYear)
		if err != nil {
			return fmt.Errorf("error retrieving events from ICS calendar: %v", err)
		}

		if len(events) == 0 {
			log.Printf("No events found in calendar %s\n", ce.Url)
			continue
		}

		for _, event := range events {
			//log.Printf("Inserting event %+v into calendar %s\n", event, calendarId)
			// Zero out the event ID to avoid invalid resource errors, let the client generate one for us
			event.Id = ""
			evt, err := calendarService.Events.Insert(calendarId, event).Do()
			if err != nil {
				return fmt.Errorf("error inserting event into calendar: %v", err)
			}
			log.Printf("Inserted event %s into calendar %s\n", evt.Id, calendarId)
		}
	}

	return nil
}

func getICSCalendarEvents(url string, start, end time.Time) ([]*calendar.Event, error) {
	log.Printf("Retrieving events from calendar: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the ICS file using gocal
	// Define a parsing window (e.g., 1 year back to 2 years forward) to capture relevant events
	c := gocal.NewParser(resp.Body)
	c.Start, c.End = &start, &end

	if err := c.Parse(); err != nil {
		return nil, err
	}

	var events []*calendar.Event
	for _, e := range c.Events {
		// Convert gocal event to Google Calendar API Event
		event := &calendar.Event{
			Summary:     e.Summary,
			Location:    e.Location,
			Description: e.Description,
			Start: &calendar.EventDateTime{
				DateTime: e.Start.Format(time.RFC3339),
			},
			End: &calendar.EventDateTime{
				DateTime: e.End.Format(time.RFC3339),
			},
			Id: e.Uid, // Use ICS UID to help with identification
		}
		events = append(events, event)
	}

	return events, nil
}

func SyncSharedCalendar(calendarService *calendar.Service, calendarId string, calendarEntries []*CalendarEntry) error {
	err := clearCalendar(calendarService, calendarId)
	if err != nil {
		return fmt.Errorf("error clearing shared calendar: %v", err)
	}

	log.Printf("Shared calendar ID: %s\n", calendarId)
	err = ImportEventsToSharedCalendar(calendarService, calendarId, calendarEntries)
	if err != nil {
		return fmt.Errorf("error importing calendar events: %v", err)
	}

	return nil
}
