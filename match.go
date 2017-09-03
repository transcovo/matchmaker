package main

import (
	"io/ioutil"
	"github.com/transcovo/matchmaker/match"
	"github.com/transcovo/matchmaker/util"
)

func main() {
	yml, err := ioutil.ReadFile("./problem.yml")
	util.PanicOnError(err, "Can't yml problem description")
	problem, err := match.LoadProblem(yml)
	match.Solve(problem)
}
