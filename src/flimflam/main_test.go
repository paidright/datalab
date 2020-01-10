package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOxford(t *testing.T) {
	result := strings.Builder{}

	assert.Nil(t, flimflam(strings.NewReader("one,two,three"), &result))

	assert.Equal(t, "one:STRING,two:STRING,three:STRING\n", result.String())
}
