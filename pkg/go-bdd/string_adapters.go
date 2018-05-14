package go_bdd

import (
	"testing"
	"reflect"
	"strings"
	"runtime"
	"fmt"
	"github.com/iancoleman/strcase"
)

func toStageFunctionArguments(test *testing.T, args ...interface{}) []reflect.Value {
	passedArguments := make([]reflect.Value, 0)
	passedArguments = append(passedArguments, reflect.ValueOf(test))
	for _, arg := range args {
		passedArguments = append(passedArguments, reflect.ValueOf(arg))
	}
	return passedArguments
}

func functionName(f reflect.Value) string {
	funcName := runtime.FuncForPC(f.Pointer()).Name()

	pos := strings.LastIndex(funcName, ".")
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + 1
	if adjustedPos >= len(funcName) {
		return ""
	}
	return funcName[adjustedPos:]
}

func toReportStringArguments(args ...interface{}) []string {
	var argumentsInString []string
	for _, argument := range args {
		argumentsInString = append(argumentsInString, fmt.Sprintf("%v", argument))
	}
	return argumentsInString
}

func toReadableSentence(functionName string, args ...interface{}) string {
	variables := toReportStringArguments(args...)
	sentence := strcase.ToSnake(functionName)
	if strings.HasPrefix(sentence, "x_") {
		sentence = "_" + sentence
	}

	for _, variable := range variables {
		sentence = strings.Replace(sentence, "_x_", fmt.Sprintf("_%v_", variable), -1)
	}
	sentence = strings.TrimSpace(strings.Replace(sentence, "_", " ", -1))

	return sentence
}
