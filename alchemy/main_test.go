package alchemy

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
)

const (
	TestNumberOfBits = 10000
	TestSetSize = 1000000  // 1 million
)


var (
	testCauldron Cauldron
)

func TestMain(m *testing.M) {
	setup()
	runExitCode := m.Run()
	teardown()

	// runExitCode 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if runExitCode == 0 && testing.CoverMode() != "" {

		coverageResult := testing.Coverage()

		// If we are less than 90% then fail the build
		if coverageResult < 0.9 {
			fmt.Printf("Tests passed but coverage failed: MUST BE >= 90%%, was %.2f\n", coverageResult*100)
			runExitCode = -1
		}
	}

	os.Exit(runExitCode)
}

func setup() {
	testCauldron = NewBitCauldron(TestSetSize)
	for k := 0; k < TestSetSize; k++ {
		testCauldron.Upsert(Item(uuid.New().String()))
	}
}

func teardown() {
	// Do something here.
}