package go_bdd

import (
	"testing"
	"reflect"
	"strings"
)

type Feature struct {
	description        string
	scenarios          []*Scenario
	beforeScenarioHook []interface{}
}

func (feature *Feature) BeforeScenario(f interface{}) {
	feature.beforeScenarioHook = append(feature.beforeScenarioHook, f)
}

type Scenario struct {
	test        *testing.T
	description string
	steps       []*Step
}

var currentFeature *Feature

func NewFeature(description string) *Feature {
	currentFeature = &Feature{
		description:        description,
		scenarios:          make([]*Scenario, 0),
		beforeScenarioHook: make([]interface{}, 0),
	}
	return currentFeature
}

func ReportFeatureTo(reporter Reporter) {
	reporter.ReportFeature(currentFeature)
}

func TestThat(t *testing.T) *Scenario {
	if currentFeature == nil {
		NewFeature("Feature Name")
	}
	description := strings.TrimLeft(t.Name(), "Test")
	description = toReadableSentence(description)
	scenario := &Scenario{
		test:        t,
		description: description,
	}

	currentFeature.scenarios = append(currentFeature.scenarios, scenario)

	for _, f := range currentFeature.beforeScenarioHook {
		reflect.ValueOf(f).Call([]reflect.Value{})
	}
	return scenario
}

func (s *Scenario) Given(f interface{}, args ...interface{}) *Scenario {
	functionToCall := reflect.ValueOf(f)
	functionArguments := toStageFunctionArguments(s.test, args...)
	functionToCall.Call(functionArguments)
	s.addStep(Given, fromTest(s.test), toReadableSentence(functionName(functionToCall), args...))
	return s
}

func (s *Scenario) When(f interface{}, args ...interface{}) *Scenario {
	functionToCall := reflect.ValueOf(f)
	functionArguments := toStageFunctionArguments(s.test, args...)
	functionToCall.Call(functionArguments)
	s.addStep(When, fromTest(s.test), toReadableSentence(functionName(functionToCall), args...))
	return s
}

func (s *Scenario) And(f interface{}, args ...interface{}) *Scenario {
	functionToCall := reflect.ValueOf(f)
	functionArguments := toStageFunctionArguments(s.test, args...)
	functionToCall.Call(functionArguments)
	s.addStep(And, fromTest(s.test), toReadableSentence(functionName(functionToCall), args...))
	return s
}

func (s *Scenario) Then(f interface{}, args ...interface{}) *Scenario {
	functionToCall := reflect.ValueOf(f)
	functionArguments := toStageFunctionArguments(s.test, args...)
	functionToCall.Call(functionArguments)
	s.addStep(Then, fromTest(s.test), toReadableSentence(functionName(functionToCall), args...))
	return s
}

func (s *Scenario) addStep(stepType stepType, stepResult stepResult, description string) {
	step := &Step{
		stepType:    stepType,
		description: description,
		result:      stepResult,
	}
	s.steps = append(s.steps, step)
}

type Step struct {
	description string
	result      stepResult
	stepType    stepType
}

var enums []string

type stepType int

func statementTypeString(s string) stepType {
	enums = append(enums, s)
	return stepType(len(enums) - 1)
}

var (
	Given = statementTypeString("Given")
	When  = statementTypeString("When")
	Then  = statementTypeString("Then")
	And   = statementTypeString("And")

	Pass = statementResultString("✓")
	Fail = statementResultString("✗")
)

func (e stepType) String() string {
	return enums[int(e)]
}

type stepResult int

func statementResultString(s string) stepResult {
	enums = append(enums, s)
	return stepResult(len(enums) - 1)
}

func fromTest(t *testing.T) stepResult {
	if t.Failed() {
		return Fail
	} else {
		return Pass
	}
}

func (e stepResult) String() string {
	return enums[int(e)]
}
