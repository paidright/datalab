package main

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrogdor(t *testing.T) {
	input := strings.NewReader(`foo,bar,baz
1,2,3
4,5,6`)
	result := strings.Builder{}
	output := csv.NewWriter(&result)
	err := trogdor(input, "baz", output)
	assert.Nil(t, err)
}
