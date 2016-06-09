package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFileInArg(t *testing.T) {
	os.Args = []string{"main", "r9w92W2Cn7MTtAhuCP5si2LH356r8FrjV.txt"}
	createArgFileIfNotExists()
	_, err := os.Stat(os.Args[1])
	assert.Equal(t, false, os.IsNotExist(err))
	os.Remove("r9w92W2Cn7MTtAhuCP5si2LH356r8FrjV.txt")
}
