package main

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	input        string
	matchOn      []matchSet
	want         []string
	demandLength int
}

func TestGumption(t *testing.T) {
	tests := []test{
		{
			input: `id,start,end
one,9am,11am
one,11am,5pm
`,
			matchOn: []matchSet{
				matchSet{
					Left:  "id",
					Right: "id",
				},
				matchSet{
					Left:  "end",
					Right: "start",
				},
			},
			want: []string{
				"one,9am,5pm,true",
			},
		},
		{
			input: `id,start,end
one,9am,11am
one,11am,5pm
two,9am,11am
two,11am,5pm
`,
			matchOn: []matchSet{
				matchSet{
					Left:  "id",
					Right: "id",
				},
				matchSet{
					Left:  "end",
					Right: "start",
				},
			},
			want: []string{
				"one,9am,5pm,true",
				"two,9am,5pm,true",
			},
		},
		{
			input: `id,start,end
one,9am,11am
one,11am,2pm
one,2pm,5pm
`,
			matchOn: []matchSet{
				matchSet{
					Left:  "id",
					Right: "id",
				},
				matchSet{
					Left:  "end",
					Right: "start",
				},
			},
			want: []string{
				"one,9am,5pm,true",
			},
		},
		{
			input: `id,start,end
one,9am,11am
one,11am,2pm
beep,bonk,bork
one,2pm,5pm
`,
			matchOn: []matchSet{
				matchSet{
					Left:  "id",
					Right: "id",
				},
				matchSet{
					Left:  "end",
					Right: "start",
				},
			},
			want: []string{
				"one,9am,2pm,true",
				"beep,bonk,bork,false",
				"one,2pm,5pm,false",
			},
		},
		{
			input: `id,paycode,start,end
one,foo,9am,11am
one,bar,11am,5pm
one,baz,9am,11am
one,quux,11am,5pm
`,
			matchOn: []matchSet{
				matchSet{
					Left:  "id",
					Right: "id",
				},
				matchSet{
					Left:  "end",
					Right: "start",
				},
				matchSet{
					LiteralRight: true,
					Left:         "paycode",
					Right:        "bar",
				},
			},
			want: []string{
				"one,bar,9am,5pm,true",
				"one,baz,9am,11am,false",
				"one,quux,11am,5pm,false",
			},
		},
		{
			input: `id,paycode,start,end
one,foo,9am,11am
one,bar,11am,5pm
one,baz,9am,11am
one,quux,11am,5pm
`,
			matchOn: []matchSet{
				matchSet{
					Left:  "id",
					Right: "id",
				},
				matchSet{
					Left:  "end",
					Right: "start",
				},
				matchSet{
					LiteralLeft: true,
					Left:        "paycode",
					Right:       "foo",
				},
			},
			want: []string{
				"one,bar,9am,5pm,true",
				"one,baz,9am,11am,false",
				"one,quux,11am,5pm,false",
			},
		},
	}

	for _, tc := range tests {
		result := strings.Builder{}

		writer := csv.NewWriter(&result)

		assert.Nil(t, ducky(strings.NewReader(tc.input), *writer, tc.matchOn))

		writer.Flush()

		for _, ex := range tc.want {
			assert.Contains(t, result.String(), ex)
		}
		if tc.demandLength != 0 {
			assert.Equal(t, tc.demandLength, len(strings.Split(result.String(), "\n")))
		}
	}
}
