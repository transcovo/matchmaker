package main

import (
	"io/ioutil"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"github.com/transcovo/go-chpr-logger"
	"time"
	"fmt"
	"github.com/transcovo/matchmaker/gcalendar"
)

func panicOnError(err error, message string) {
	if err != nil {
		logger.GetLogger().WithError(err).Fatal(message)
	}
}

func main() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	panicOnError(err, "Can't read credentials file")

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	panicOnError(err, "Can't load client configuration")

	client := gcalendar.GetClient(ctx, config)

	srv, err := calendar.New(client)
	panicOnError(err, "Can't retrieve calendar client")

	event, err := srv.Events.Insert("chauffeur-prive.com_k23ttdrv7g0l5i2vjj1f3s8voc@group.calendar.google.com", &calendar.Event{
		Start:&calendar.EventDateTime{
			DateTime:time.Now().Format(time.RFC3339),
		},
		End:&calendar.EventDateTime{
			DateTime:time.Now().Add(time.Hour).Format(time.RFC3339),
		},
		Summary:"Test API",
		Attendees:[]*calendar.EventAttendee{
			{
				Email:"samuel@chauffeur-prive.com",
			},
		},
	}).Do()
	panicOnError(err, "Can't create event")

	fmt.Printf("Event created: %s\n", event.Summary)

	/*
	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("samuel.rossille@gmail.com").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	panicOnError(err, "Can't events")

	fmt.Println("Upcoming events:")
	if len(events.Items) > 0 {
		for _, i := range events.Items {
			var when string
			// If the DateTime is an empty string the Event is an all-day Event.
			// So only Date is available.
			if i.Start.DateTime != "" {
				when = i.Start.DateTime
			} else {
				when = i.Start.Date
			}
			fmt.Printf("  - %s (%s)\n", i.Summary, when)
		}
	} else {
		fmt.Printf("No upcoming events found.\n")
	}*/
}
