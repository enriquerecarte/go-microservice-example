package go_bdd

import (
	"fmt"
	"strings"
)

type Reporter interface {
	ReportFeature(*Feature)
	Flush()
}

type ConsoleReporter struct {
	report string
}

const FeatureLayout = "Feature: %s"
const ScenarioLayout = "Scenario: %s"
const StepLayout = "%s %s %s"

func (r *ConsoleReporter) ReportFeature(feature *Feature) {
	r.writeLine(fmt.Sprintf(FeatureLayout, feature.description))
	for _, scenario := range feature.scenarios {
		r.writeLine("")
		r.writeLine("\t" + fmt.Sprintf(ScenarioLayout, strings.Title(scenario.description)))

		for _, step := range scenario.steps {
			r.writeLine("\t\t" + fmt.Sprintf(StepLayout, step.stepType.String(), step.description, step.result.String()))
		}
	}
}

func (r *ConsoleReporter) Flush() {
	fmt.Println(r.report)
}

func (r *ConsoleReporter) writeLine(line string) {
	if r.report != "" {
		r.report += "\n"
	}
	r.report += line
}