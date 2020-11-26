package alchemy

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

const (
	TestNumberOfBits = 10000
	TestSetMaxSize = 1000000
)

var (
	testCauldron Cauldron
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	testCauldron = newBitCauldron()
	for k := 0; k < TestSetMaxSize; k++ {
		testCauldron.Upsert(Item(uuid.New().String()))
	}
}

func teardown() {
	// Do something here.
}