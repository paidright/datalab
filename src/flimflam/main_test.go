package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlimFlam(t *testing.T) {
	result := strings.Builder{}

	assert.Nil(t, flimflam(strings.NewReader("one,two,three"), "kv", &result))

	assert.Equal(t, "one:STRING,two:STRING,three:STRING\n", result.String())
}

func TestJsonOutput(t *testing.T) {
	result := strings.Builder{}

	assert.Nil(t, flimflam(strings.NewReader("one,two,three"), "json", &result))

	assert.Equal(t, `[{"name":"one","type":"STRING"},{"name":"two","type":"STRING"},{"name":"three","type":"STRING"}]`, result.String())
}
