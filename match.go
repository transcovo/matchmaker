package main

import (
	"io/ioutil"
	"github.com/rossille/matchmaker/match"
	"github.com/rossille/matchmaker/util"
	"flag"
	"os"
	"log"
	"runtime/pprof"
	"fmt"
	"github.com/rossille/matchmaker/gcalendar"
	"google.golang.org/api/calendar/v3"
	"github.com/transcovo/go-chpr-logger"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	yml, err := ioutil.ReadFile("./problem.yml")
	util.PanicOnError(err, "Can't yml problem description")
	problem, err := match.LoadProblem(yml)
	solution := match.Solve(problem)

	var response string

	cal, err := gcalendar.GetGoogleCalendarService()
	util.PanicOnError(err, "Can't get gcalendar client")

	for _, session := range solution.Sessions {
		fmt.Scanln(&response)

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
