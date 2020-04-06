# Matchmaker

Matchmaker takes care of matching and planning of reviewers and review slots in people's calendars.

## Setup

You need to retrieve the `persons.yml` file containing people configuration for review.
Format example:
```yaml
- email: john.doe@kapten.com
  isgoodreviewer: true
- email: chuck.norris@kapten.com
  isgoodreviewer: true
  maxsessionsperweek: 1
- email: james.bond@kapten.com
- email: john.wick@kapten.com
  maxsessionsperweek: 1
- email: obi-wan.kenobi@kapten.com
```
**isgoodreviewer** is used to distinguish the experienced reviewers in order to create reviewer pairs that contain at least one experienced reviewer.
**maxsessionsperweek** is optional, default is 3. If set to 0, it also falls back to the default value.

You need to create/retrieve a `client_secret.json` file containing a valid Google Calendar
access token for Kapten's calendar.

Those files need to be placed at the root of the project.

## Preparing

    go run prepare.go [-week-shift value [default=0]]

This script will compute work ranges for the target week, and check free slots for each potential
reviewer and create an output file `problem.yml`.

By default, the script plans for the upcoming monday, you can provide a `weekShift` value as a parameter, allowing
to plan for further weeks (1 = the week after upcoming monday, etc.)

## Matching

    go run match.go

This script will take input from the `match.yml` file and match reviewers together in review slots for the target week.
The output is a `planning.yml` file with reviewers couples and planned slots.

## Planning

    go run plan.go

This script will take input from the `planning.yml` file and create review events in reviewers' calendar.


## Default run

By running the script:
    ./do.sh

All scripts will be run sequentially for the upcoming monday.
