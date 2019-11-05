package main

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessFile(t *testing.T) {
	input := strings.NewReader(`hi,foo,bar
1,2,3
a,b,c`)

	targets := []target{
		target{
			sources: []string{"foo", "bar"},
			sep:     "-",
			dest:    "baz",
		},
	}

	result := strings.Builder{}

	writer := csv.NewWriter(&result)

	assert.Nil(t, processFile(input, targets, *writer))

	writer.Flush()

	for _, ex := range []string{
		"hi,foo,bar,baz", "a,b,c,b-c", "1,2,3,2-3",
	} {
		assert.Contains(t, result.String(), ex)
	}
}
