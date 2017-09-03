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
	util.PanicOnError(err, "Can't parse yml problem description")
	ymlOut, _ := problem.ToYaml()
	println(string(ymlOut))
}
