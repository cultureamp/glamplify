package alchemy

import (
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
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	testCauldron = newBitCauldron(TestSetSize)
	for k := 0; k < TestSetSize; k++ {
		testCauldron.Upsert(Item(uuid.New().String()))
	}
}

func teardown() {
	// Do something here.
}