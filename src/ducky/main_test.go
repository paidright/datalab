package main

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	name         string
	input        string
	matchOn      []matchSet
	want         []string
	demandLength int
}

func TestGumption(t *testing.T) {
	tests := []test{
		{
			name: "basic",
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
			name: "multiline basic",
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
			name: "three way merge",
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
			name: "interruption",
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
			name: "literal right",
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
			name: "literal left",
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
		{
			name: "non contiguous",
			input: `id,start,end,ducky_taped
one,9am,11am,true
one,11am,5pm,false
two,9am,10am,false
two,11am,5pm,false
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
				"two,9am,10am,false",
				"two,11am,5pm,false",
			},
		},
		{
			name: "inverse literal right",
			input: `id,start,end
one,9am,11am
one,11am,5pm
two,9am,11am
two,11am,never
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
					Inverse:      true,
					LiteralRight: true,
					Left:         "end",
					Right:        "never",
				},
			},
			want: []string{
				"one,9am,5pm,true",
				"two,9am,11am,false",
				"two,11am,never,false",
			},
		},
		{
			name: "inverse literal right",
			input: `id,start,end,classification
one,9am,11am,foo
one,11am,wut,foo
one,11am,wut,foo
`,
			matchOn: []matchSet{
				matchSet{
					Left:  "id",
					Right: "id",
				},
				matchSet{
					Left:  "classification",
					Right: "classification",
				},
				matchSet{
					Left:  "end",
					Right: "start",
				},
				matchSet{
					Inverse:      true,
					LiteralRight: true,
					Left:         "end",
					Right:        "wut",
				},
			},
			want: []string{
				"one,9am,11am,foo,false",
				"one,11am,wut,foo,false",
				"one,11am,wut,foo,false",
			},
		},
		{
			name: "double inverse literal right",
			input: `id,start,end
one,9am,11am
one,11am,5pm
two,never,never
two,never,never
`,
			matchOn: []matchSet{
				matchSet{
					Left:  "end",
					Right: "start",
				},
				matchSet{
					Left:  "id",
					Right: "id",
				},
				matchSet{
					Inverse:      true,
					LiteralRight: true,
					Left:         "end",
					Right:        "never",
				},
				matchSet{
					Inverse:      true,
					LiteralRight: true,
					Left:         "start",
					Right:        "never",
				},
			},
			want: []string{
				"one,9am,5pm,true",
				"two,never,never,false",
				"two,never,never,false",
			},
		},
		{
			name: "inverse basic match",
			input: `id,start,end
one,9am,11am
one,12am,5pm
two,9am,11am
two,11am,5pm
`,
			matchOn: []matchSet{
				matchSet{
					Left:  "id",
					Right: "id",
				},
				matchSet{
					Inverse: true,
					Left:    "end",
					Right:   "start",
				},
			},
			want: []string{
				"one,9am,5pm,true",
				"two,9am,11am,false",
				"two,11am,5pm,false",
			},
		},
	}

	for _, tc := range tests {
		result := strings.Builder{}

		writer := csv.NewWriter(&result)

		assert.Nil(t, ducky(strings.NewReader(tc.input), *writer, tc.matchOn))

		writer.Flush()

		for _, ex := range tc.want {
			assert.Contains(t, result.String(), ex, tc.name)
		}
		if tc.demandLength != 0 {
			assert.Equal(t, tc.demandLength, len(strings.Split(result.String(), "\n")))
		}
	}
}
