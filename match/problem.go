package match

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Level int

const (
	LanguageLevelNone       Level = iota
	LanguageLevelApprentice
	LanguageLevelMaster
)

type Exclusivity int

func (exclusivity Exclusivity) String() string {
	switch exclusivity {
	case ExclusivityNone:
		return "none"
	case ExclusivityBack:
		return "back"
	case ExclusivityMobile:
		return "mobile"
	}
	panic("Unexpected exclusivity " + string(exclusivity))
}

const (
	ExclusivityNone Exclusivity = iota
	ExclusivityMobile
	ExclusivityBack
)

type Languages struct {
	Js    Level
	Go  Level
	Python  Level
	Ios     Level
	Android Level
}

func (languages *Languages) GetExclusivity() Exclusivity {
	if languages.Js == LanguageLevelNone && languages.Go == LanguageLevelNone && languages.Python == LanguageLevelNone {
		return ExclusivityMobile
	}
	if languages.Ios == LanguageLevelNone && languages.Android == LanguageLevelNone {
		return ExclusivityBack
	}
	return ExclusivityNone
}

type Person struct {
	Email          string `yaml:"email"`
	Languages      Languages
	IsGoodReviewer bool
}

func (person *Person) GetExclusivity() Exclusivity {
	return person.Languages.GetExclusivity()
}

func LoadPersons(path string) ([]*Person, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var persons []*Person
	yaml.Unmarshal(data, &persons)

	return persons, nil
}

type BusyTime struct {
	Person *Person
	Range  *Range
}

type Problem struct {
	People         []*Person
	WorkRanges     []*Range
	BusyTimes      []*BusyTime
	TargetCoverage map[Exclusivity]int
}

type SerializedBusyTime struct {
	Email string
	Range *Range
}

type SerializedProblem struct {
	People         []*Person
	WorkRanges     []*Range
	BusyTimes      []*SerializedBusyTime
	TargetCoverage map[Exclusivity]int
}

func (problem *Problem) ToYaml() ([]byte, error) {
	serializedBusyTimes := make([]*SerializedBusyTime, len(problem.BusyTimes))
	for i, busyTime := range problem.BusyTimes {
		serializedBusyTimes[i] = &SerializedBusyTime{
			Email: busyTime.Person.Email,
			Range: busyTime.Range,
		}
	}

	serializedProblem := SerializedProblem{
		People:         problem.People,
		WorkRanges:     problem.WorkRanges,
		BusyTimes:      serializedBusyTimes,
		TargetCoverage: problem.TargetCoverage,
	}
	data, err := yaml.Marshal(serializedProblem)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func LoadProblem(yml []byte) (*Problem, error) {
	var serializedProblem SerializedProblem
	err := yaml.Unmarshal(yml, &serializedProblem)
	if err != nil {
		return nil, err
	}

	personsByEmail := map[string]*Person{}
	for _, person := range serializedProblem.People {
		personsByEmail[person.Email] = person
	}

	busyTimes := make([]*BusyTime, len(serializedProblem.BusyTimes))
	for i, serializedBusyTime := range serializedProblem.BusyTimes {
		busyTimes[i] = &BusyTime{
			Person: personsByEmail[serializedBusyTime.Email],
			Range:  serializedBusyTime.Range,
		}
	}

	return &Problem{
		People:         serializedProblem.People,
		WorkRanges:     serializedProblem.WorkRanges,
		BusyTimes:      busyTimes,
		TargetCoverage: serializedProblem.TargetCoverage,
	}, nil
}
