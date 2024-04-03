package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestValidateFilename(t *testing.T) {
	err := validateFilename("test.dump")
	assert.NoError(t, err)
}

func TestCleanup(t *testing.T) {
	file, err := os.Create("./test.dump")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	os.Setenv("FILENAME", "test.dump")
	cleanup()
	_, err = os.Stat("test.dump")
	// File exists and cleanup has not done its job
	assert.Error(t, err)
}
