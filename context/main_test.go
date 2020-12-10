package context

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	runExitCode := m.Run()

	// runExitCode 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if runExitCode == 0 && testing.CoverMode() != "" {

		coverageResult := testing.Coverage()

		// If we are less than 80% then fail the build
		if coverageResult < 0.8 {
			fmt.Printf("Tests passed but coverage failed: MUST BE >= 80%%, was %.2f\n", coverageResult*100)
			runExitCode = -1
		}
	}
	os.Exit(runExitCode)
}

