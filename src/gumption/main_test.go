package main

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	input string
	flags map[string]flagval
	cols  []string
	want  []string
}

func TestGumption(t *testing.T) {
	tests := []test{
		{
			flags: map[string]flagval{
				"commasToPoints": flagval{
					active: true,
				},
			},
			input: `one,two
"123,456",abc`,
			want: []string{"123.456,abc"},
		},
		{
			flags: map[string]flagval{
				"stripLeadingZeroes": flagval{
					active: true,
				},
			},
			input: `one,two
0000123,00abc`,
			want: []string{"123,abc"},
		},
		{
			flags: map[string]flagval{
				"unquote": flagval{
					active: true,
				},
			},
			input: `one,two
'123',abc`,
			want: []string{"123,abc"},
		},
		{
			flags: map[string]flagval{
				"addMissing": flagval{
					active: true,
					value:  "asd",
				},
			},
			input: `one,two
,abc`,
			want: []string{"asd,abc"},
		},
		{
			flags: map[string]flagval{
				"replaceCell": flagval{
					active: true,
					replacements: []replacement{
						{
							from: "123",
							to:   "xyz",
						},
						{
							from: "456",
							to:   "qwe",
						},
					},
				},
			},
			input: `one,two
456,poi
123,abc`,
			want: []string{"xyz,abc", "qwe,poi"},
		},
		{
			flags: map[string]flagval{
				"rename": flagval{
					active: true,
					value:  "asd",
				},
			},
			cols: []string{"two"},
			input: `one,two
123,abc`,
			want: []string{"one,asd"},
		},
		{
			flags: map[string]flagval{
				"splitOnDelim": flagval{
					active: true,
					value:  "-",
				},
			},
			cols: []string{"one"},
			input: `one,two
123-456,abc`,
			want: []string{"one,two,one_1", "123,abc,456"},
		},
		{
			flags: map[string]flagval{
				"cp": flagval{
					active: true,
				},
			},
			cols: []string{"one"},
			input: `one,two
123,abc`,
			want: []string{"one,two,one_1", "123,abc,123"},
		},
		{
			flags: map[string]flagval{
				"drop": flagval{
					active: true,
				},
			},
			cols: []string{"three"},
			input: `one,two,three
123,abc,3`,
			want: []string{"one,two", "123,abc"},
		},
		{
			flags: map[string]flagval{
				"stripLeadingZeroes": flagval{
					active: true,
				},
				"commasToPoints": flagval{
					active: true,
				},
				"splitOnDelim": flagval{
					active: true,
					value:  "-",
				},
				"addMissing": flagval{
					active: true,
					value:  "999",
				},
			},
			cols: []string{"one"},
			input: `one,two,three
"0001,23-456",abc,3
,abc,3`,
			want: []string{
				`one,two,three,one_1`,
				`999,abc,3,`,
				`1.23,abc,3,456`,
			},
		},
	}

	for _, tc := range tests {
		result := strings.Builder{}

		writer := csv.NewWriter(&result)

		assert.Nil(t, gumption(strings.NewReader(tc.input), *writer, tc.cols, tc.flags))

		writer.Flush()

		for _, ex := range tc.want {
			assert.Contains(t, result.String(), ex)
		}
	}
}

func TestParseReplacements(t *testing.T) {
	expected := []replacement{
		{
			from: "A",
			to:   "B",
		}, {
			from: "X",
			to:   "Y",
		},
	}

	input := "A,B,X,Y"

	assert.Equal(t, expected, parseReplacements(input))
}
