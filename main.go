package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/t-hale/family-cal/lib"
	"google.golang.org/api/calendar/v3"
)

var (
	calendarFile = flag.String("calendar-file", "", "Path to the file containing calendar URLs")
	calendarName = flag.String("calendar-name", "", "Name of the target shared calendar")
)

func main() {

	flag.Parse()

	if *calendarFile == "" || *calendarName == "" {
		flag.Usage()
		log.Fatalf("Both --calendar-file and --calendar-name flags are required.")
	}

	log.Printf("Reading calendar URLs from file: %s\n", *calendarFile)
	log.Printf("Target calendar: %s\n", *calendarName)

	ctx := context.Background()
	calendarService, err := calendar.NewService(ctx)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	sharedCalendarId, err := lib.CreateSharedCalendar(calendarService, "RVA Hale Sports")
	if err != nil {
		log.Fatalf("Error creating shared calendar: %v", err)
	}

	calendarEntries, err := readCalendarFile()
	if err != nil {
		log.Fatalf("Error reading calendar file: %v", err)
	}

	err = lib.SyncSharedCalendar(calendarService, sharedCalendarId, calendarEntries)
	if err != nil {
		log.Fatalf("Error syncing shared calendar: %v", err)
	}
}

func readCalendarFile() ([]*lib.CalendarEntry, error) {
	file, err := os.Open(*calendarFile)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	firstLine := true
	calendarEntries := []*lib.CalendarEntry{}
	for scanner.Scan() {
		line := scanner.Text()
		if firstLine {
			firstLine = false
			continue
		}
		if line == "" {
			continue
		}

		items := strings.Split(line, ",")
		ce := &lib.CalendarEntry{
			Name: items[0],
			Url:  items[1],
		}

		calendarEntries = append(calendarEntries, ce)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return calendarEntries, nil
}
