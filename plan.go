package main

import (
	"github.com/transcovo/go-chpr-logger"
	"github.com/transcovo/matchmaker/gcalendar"
	"github.com/transcovo/matchmaker/match"
	"github.com/transcovo/matchmaker/util"
	"google.golang.org/api/calendar/v3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func LoadPlan(yml []byte) (*match.Solution, error) {
	var solution *match.Solution
	err := yaml.Unmarshal(yml, &solution)
	if err != nil {
		return nil, err
	}

	return solution, nil
}

func main() {
	yml, err := ioutil.ReadFile("./planning.yml")
	util.PanicOnError(err, "Can't yml problem description")

	cal, err := gcalendar.GetGoogleCalendarService()
	util.PanicOnError(err, "Can't get gcalendar client")

	solution, err := LoadPlan(yml)
	util.PanicOnError(err, "Can't get solution from planning file")

	for _, session := range solution.Sessions {
		attendees := []*calendar.EventAttendee{}

		for _, person := range session.Reviewers.People {
			attendees = append(attendees, &calendar.EventAttendee{
				Email: person.Email,
			})
		}

		_, err := cal.Events.Insert("chauffeur-prive.com_k23ttdrv7g0l5i2vjj1f3s8voc@group.calendar.google.com", &calendar.Event{
			Start: &calendar.EventDateTime{
				DateTime: gcalendar.FormatTime(session.Range.Start),
			},
			End: &calendar.EventDateTime{
				DateTime: gcalendar.FormatTime(session.Range.End),
			},
			Summary: session.GetDisplayName(),
			Attendees: attendees,
		}).Do()
		util.PanicOnError(err, "Can't create event")
		logger.Info("✔ " + session.GetDisplayName())
	}
}
