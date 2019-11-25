package main

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	delim rune
	input string
	want  []string
}

func TestOxford(t *testing.T) {
	tests := []test{
		{
			delim: '|',
			input: `one|two
oh|hai`,
			want: []string{"one,two", "oh,hai"},
		},
		{
			delim: '|',
			input: `one|amount
oh|3,234.10`,
			want: []string{"one,amount", `oh,"3,234.10"`},
		},
	}

	for _, tc := range tests {
		result := strings.Builder{}

		writer := csv.NewWriter(&result)

		assert.Nil(t, oxford(strings.NewReader(tc.input), writer, tc.delim))

		writer.Flush()

		for _, ex := range tc.want {
			assert.Contains(t, result.String(), ex)
		}
	}
}
