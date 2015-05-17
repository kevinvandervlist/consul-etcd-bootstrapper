package util

import (
	"log"
	"testing"
)

func AssertNoErrorCallback(error error, f func(err error)) {
	if error != nil {
		log.Printf("Error occurred: %s\n", error)
		f(error)
	}
}

func AssertNoError(err error) {
	AssertNoErrorCallback(err, func(e error) {})
}

func AssertEquals(expected interface{}, actual interface{}, t *testing.T) {
	if expected != actual {
		t.Fatalf("Expected '%v', got '%v'.\n", expected, actual)
	}
}
