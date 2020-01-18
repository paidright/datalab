package main

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	input        string
	flags        map[string]flagval
	cols         []string
	want         []string
	demandLength int
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
				"replaceCellLookup": flagval{
					active: true,
					replacements: []replacement{
						{
							from: "123",
							to:   "two",
						},
						{
							from: "456",
							to:   "two",
						},
					},
				},
			},
			input: `one,two
456,poi
123,abc
789,xyz`,
			want: []string{"poi,poi", "abc,abc", "789,xyz"},
		},
		{
			flags: map[string]flagval{
				"replaceChar": flagval{
					active: true,
					replacements: []replacement{
						{
							from: ":",
							to:   ".",
						},
						{
							from: ";",
							to:   ".",
						},
					},
				},
			},
			cols: []string{"one"},
			input: `one,two
4:56,7:89
1;23,abc`,
			want: []string{"4.56,7:89", "1.23,abc"},
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
				"rename": flagval{
					active: true,
					value:  "asd",
				},
			},
			cols: []string{"somethingGUMPTION_LITERAL_COMMAelse"},
			input: `"something,else",two
123,abc`,
			want: []string{"asd,two"},
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
				"stompAlphas": flagval{
					active: true,
				},
			},
			cols: []string{"two"},
			input: `one,two,three
123,1,abc
123,a,xyz
123,ab2a,abc`,
			want: []string{"one,two,three", "123,1,abc", "123,,xyz", "123,2,abc"},
		},
		{
			flags: map[string]flagval{
				"deleteWhere": flagval{
					active: true,
					value:  "xyz",
				},
			},
			cols: []string{"three"},
			input: `one,two,three
123,1,abc
123,a,xyz
123,2,abc`,
			want:         []string{"one,two,three", "123,1,abc", "123,2,abc"},
			demandLength: 4,
		},
		{
			flags: map[string]flagval{
				"deleteWhereNot": flagval{
					active: true,
					value:  "abc",
				},
			},
			cols: []string{"three"},
			input: `one,two,three
123,1,abc
123,a,xyz
123,2,abc`,
			want:         []string{"one,two,three", "123,1,abc", "123,2,abc"},
			demandLength: 4,
		},
		{
			flags: map[string]flagval{
				"trimWhitespace": flagval{
					active: true,
				},
			},
			cols: []string{"one"},
			input: `one,two
 123 ,abc
1 23,abc`,
			want: []string{"one,two", "123,abc", "1 23,abc"},
		},
		{
			flags: map[string]flagval{
				"backToFront": flagval{
					active: true,
					value:  "-",
				},
			},
			cols: []string{"one"},
			input: `one,two
123-,abc
1.23-,abc
,abc
-1.45,abc`,
			want: []string{"one,two", "-123,abc", "-1.23,abc", ",abc", "-1.45,abc"},
		},
		{
			flags: map[string]flagval{
				"reformatDate": flagval{
					active: true,
					value:  "DD.MM.YYYY,YYYY-MM-DD",
				},
			},
			cols: []string{"one"},
			input: `one,two
lolwut,hurr
21.07.2003,foo`,
			want: []string{"one,two", "lolwut,hurr", "2003-07-21,foo"},
		},
		{
			flags: map[string]flagval{
				"reformatDate": flagval{
					active: true,
					value:  "DD.SHORTMONTH.YYYY,YYYY-MM-DD",
				},
			},
			cols: []string{"one"},
			input: `one,two
lolwut,hurr
21.MAR.2003,foo`,
			want: []string{"one,two", "lolwut,hurr", "2003-03-21,foo"},
		},
		{
			flags: map[string]flagval{
				"reformatDate": flagval{
					active: true,
					value:  "DD.SHORTMONTH.YY,YYYY-MM-DD",
				},
			},
			cols: []string{"one"},
			input: `one,two
lolwut,hurr
04.NOV.16,foo
21.MAR.03,foo`,
			want: []string{"one,two", "lolwut,hurr", "2016-11-04,foo", "2003-03-21,foo"},
		},
		{
			flags: map[string]flagval{
				"cleanCols": flagval{
					active: true,
				},
			},
			cols: []string{},
			input: `with space, and whitespace  ,got.dots,maybe-a-dash,  all.together-now
lolwut,hurr,foo,bar,baz`,
			want: []string{
				"with_space,and_whitespace,got_dots,maybe_a_dash,all_together_now",
				"lolwut,hurr,foo,bar,baz",
			},
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
		{
			flags: map[string]flagval{
				"leftPad": flagval{
					active: true,
					value:  "0,4",
				},
			},
			cols:  []string{"one"},
			input: "one,two\n1,1\n11,1\n1111,11",
			want: []string{
				"one,two",
				"0001,1",
				"0011,1",
				"1111,11",
			},
		},
		{
			flags: map[string]flagval{
				"reformatTime": flagval{
					active: true,
					value:  "HHMM,HH:MM",
				},
			},
			cols: []string{"one"},
			input: `one,two
0830,foo`,
			want: []string{"one,two", "08:30,foo"},
		},
		{
			flags: map[string]flagval{
				"reformatDateTime": flagval{
					active: true,
					value:  "YYYYMMDDhhmmss,YYYY-MM-DD hh:mm:ss",
				},
			},
			cols: []string{"one"},
			input: `one,two
20150629083000,foo`,
			want: []string{"one,two", "2015-06-29 08:30:00,foo"},
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
		if tc.demandLength != 0 {
			assert.Equal(t, tc.demandLength, len(strings.Split(result.String(), "\n")))
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
