package main

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessFile(t *testing.T) {
	left := strings.NewReader(`id,foo
1,a
2,x`)
	right := strings.NewReader(`id,bar
1,b
2,y
1,0`)
	joinKey := "id"

	result := strings.Builder{}
	output := csv.NewWriter(&result)

	assert.Nil(t, join(joinKey, left, right, output))

	output.Flush()

	// Len is 5 due to a trailing newline
	assert.Equal(t, 5, len(strings.Split(result.String(), "\n")))

	expected := []string{
		"id,foo,bar,left_original_line_number,right_original_line_number",
		"1,a,0,2,4",
		"2,x,y,3,3",
		"1,a,b,2,2",
	}

	for _, line := range expected {
		assert.Contains(t, result.String(), line)
	}
}

func TestMissingRight(t *testing.T) {
	left := strings.NewReader(`id,foo
1,a
2,x`)
	right := strings.NewReader(`id,bar
1,b
2,y
4,y
1,0`)
	joinKey := "id"

	result := strings.Builder{}
	output := csv.NewWriter(&result)

	assert.Nil(t, join(joinKey, left, right, output))

	output.Flush()

	// Len is 5 due to a trailing newline
	assert.Equal(t, 6, len(strings.Split(result.String(), "\n")))

	expected := []string{
		"id,foo,bar,left_original_line_number,right_original_line_number",
		"1,a,0,2,5",
		"2,x,y,3,3",
		"4,,y,,4",
		"1,a,0,2,5",
	}

	for _, line := range expected {
		assert.Contains(t, result.String(), line)
	}
}
