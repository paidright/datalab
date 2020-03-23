package main

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/paidright/datalab/util"
	"github.com/stretchr/testify/assert"
)

type test struct {
	name         string
	input        string
	matchOn      []matchSet
	want         []string
	demandLength int
}

type matchTest struct {
	name  string
	left  util.Line
	right util.Line
	match matchSet
	want  bool
}

type anyMatchTest struct {
	name  string
	group []util.Line
	match matchSet
	want  bool
}

func TestDoLinesMatch(t *testing.T) {
	tests := []matchTest{
		{
			name: "basic match",
			left: util.Line{
				Data: map[string]string{
					"start": "11am",
					"end":   "5pm",
				},
			},
			right: util.Line{
				Data: map[string]string{
					"start": "9am",
					"end":   "11am",
				},
			},
			match: matchSet{
				Left:  "start",
				Right: "end",
			},
			want: true,
		},
		{
			name: "basic miss",
			left: util.Line{
				Data: map[string]string{
					"start": "1pm",
					"end":   "5pm",
				},
			},
			right: util.Line{
				Data: map[string]string{
					"start": "9am",
					"end":   "11am",
				},
			},
			match: matchSet{
				Left:  "start",
				Right: "end",
			},
			want: false,
		},
		{
			name: "inverse basic match",
			left: util.Line{
				Data: map[string]string{
					"start": "11am",
					"end":   "5pm",
				},
			},
			right: util.Line{
				Data: map[string]string{
					"start": "9am",
					"end":   "11am",
				},
			},
			match: matchSet{
				Left:    "start",
				Right:   "end",
				Inverse: true,
			},
			want: false,
		},
		{
			name: "inverse basic miss",
			left: util.Line{
				Data: map[string]string{
					"start": "1pm",
					"end":   "5pm",
				},
			},
			right: util.Line{
				Data: map[string]string{
					"start": "9am",
					"end":   "11am",
				},
			},
			match: matchSet{
				Left:    "start",
				Right:   "end",
				Inverse: true,
			},
			want: true,
		},
		{
			name: "literal right match",
			left: util.Line{
				Data: map[string]string{
					"paycode": "foo",
					"end":     "11am",
				},
			},
			right: util.Line{
				Data: map[string]string{
					"paycode": "bar",
					"end":     "5pm",
				},
			},
			match: matchSet{
				LiteralRight: true,
				Left:         "paycode",
				Right:        "bar",
			},
			want: true,
		},
		{
			name: "literal left match",
			left: util.Line{
				Data: map[string]string{
					"paycode": "foo",
					"end":     "11am",
				},
			},
			right: util.Line{
				Data: map[string]string{
					"paycode": "bar",
					"end":     "5pm",
				},
			},
			match: matchSet{
				LiteralLeft: true,
				Left:        "paycode",
				Right:       "foo",
			},
			want: true,
		},
		{
			name: "literal right miss",
			left: util.Line{
				Data: map[string]string{
					"paycode": "foo",
					"end":     "11am",
				},
			},
			right: util.Line{
				Data: map[string]string{
					"paycode": "bar",
					"end":     "5pm",
				},
			},
			match: matchSet{
				LiteralRight: true,
				Left:         "paycode",
				Right:        "foo",
			},
			want: false,
		},
		{
			name: "literal left miss",
			left: util.Line{
				Data: map[string]string{
					"paycode": "foo",
					"end":     "11am",
				},
			},
			right: util.Line{
				Data: map[string]string{
					"paycode": "bar",
					"end":     "5pm",
				},
			},
			match: matchSet{
				LiteralLeft: true,
				Left:        "paycode",
				Right:       "bar",
			},
			want: false,
		},
		{
			name: "inverse literal right match",
			left: util.Line{
				Data: map[string]string{
					"paycode": "foo",
					"end":     "11am",
				},
			},
			right: util.Line{
				Data: map[string]string{
					"paycode": "quux",
					"end":     "5pm",
				},
			},
			match: matchSet{
				Inverse:      true,
				LiteralRight: true,
				Left:         "paycode",
				Right:        "bar",
			},
			want: true,
		},
		{
			name: "inverse literal right miss",
			left: util.Line{
				Data: map[string]string{
					"paycode": "foo",
					"end":     "11am",
				},
			},
			right: util.Line{
				Data: map[string]string{
					"paycode": "bar",
					"end":     "5pm",
				},
			},
			match: matchSet{
				Inverse:      true,
				LiteralRight: true,
				Left:         "paycode",
				Right:        "bar",
			},
			want: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, doLinesMatch(tc.left, tc.right, tc.match), tc.name)
		})
	}
}

func TestAnyLinesMatch(t *testing.T) {
	tests := []anyMatchTest{
		{
			name: "basic match",
			group: []util.Line{
				util.Line{
					Data: map[string]string{
						"start": "11am",
						"end":   "5pm",
					},
				},
				util.Line{
					Data: map[string]string{
						"start": "9am",
						"end":   "11am",
					},
				},
			},
			match: matchSet{
				Left:  "start",
				Right: "end",
			},
			want: true,
		},
		{
			name: "literal left match",
			group: []util.Line{
				util.Line{
					Data: map[string]string{
						"paycode": "bar",
						"end":     "5pm",
					},
				},
				util.Line{
					Data: map[string]string{
						"paycode": "foo",
						"end":     "11am",
					},
				},
			},
			match: matchSet{
				LiteralLeft: true,
				Left:        "paycode",
				Right:       "foo",
			},
			want: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, anyLinesMatch(tc.group, tc.match), tc.name)
		})
	}
}

func TestDucky(t *testing.T) {
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
			name: "three way with any",
			input: `id,start,end,flag
foo,9am,11am,yep
foo,11am,2pm,nope
foo,2pm,5pm,nope
one,9am,11am,nope
one,11am,2pm,yep
one,2pm,5pm,nope
two,9am,11am,nope
two,11am,2pm,nope
two,2pm,5pm,nope
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
					MatchAny:     true,
					LiteralRight: true,
					Left:         "flag",
					Right:        "yep",
				},
			},
			want: []string{
				"foo,9am,5pm,yep,true",
				"one,9am,5pm,nope,true",
				"two,9am,11am,nope,false",
				"two,11am,2pm,nope,false",
				"two,2pm,5pm,nope,false",
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
never,never,never
never,never,never
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
					Left:         "id",
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
		t.Run(tc.name, func(t *testing.T) {
			result := strings.Builder{}

			writer := csv.NewWriter(&result)

			assert.Nil(t, ducky(strings.NewReader(tc.input), *writer, tc.matchOn, "id"))

			writer.Flush()

			for _, ex := range tc.want {
				assert.Contains(t, result.String(), ex, tc.name)
			}
			if tc.demandLength != 0 {
				assert.Equal(t, tc.demandLength, len(strings.Split(result.String(), "\n")))
			}
		})
	}
}
