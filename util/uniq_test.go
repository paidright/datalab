package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var uniqTests = []struct {
	name string
	in   []string
	out  []string
}{
	{"no dupes", []string{"one", "two"}, []string{"one", "two"}},
	{"dupes", []string{"one", "two", "one"}, []string{"one", "two"}},
	{"tripes", []string{"one", "one", "two", "one"}, []string{"one", "two"}},
	{"multidupes", []string{"one", "two", "one", "two", "one", "two", "three", "three"}, []string{"one", "two", "three"}},
}

func TestUniq(t *testing.T) {
	for _, uT := range uniqTests {
		t.Run(uT.name, func(t *testing.T) {
			s := Uniq(uT.in)
			assert.Equal(t, uT.out, s, "got %q, want %q")
		})
	}
}
