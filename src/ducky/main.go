package main

import (
	"encoding/csv"
	"flag"
	"io"
	"os"
	"strings"

	"github.com/paidright/datalab/util"
)

var version = flag.Bool("version", false, "Just print the version and exit")
var quiet = flag.Bool("quiet", false, "Tone down the output noise")
var matchInput = flag.String("match", "", "A comma separated list of columns to match on. eg: id:id,end:start")
var inverseMatchInput = flag.String("inverse-match", "", "A comma separated list of columns to inverse match on. eg: id:id,end:start")
var literalLeftMatchInput = flag.String("match-literal-left", "", "A comma separated list values to match the left branch on. eg: paycode:salary")
var literalRightMatchInput = flag.String("match-literal-right", "", "A comma separated list values to match the right branch on. eg: paycode:extra_hours")
var inverseLiteralLeftMatchInput = flag.String("inverse-match-literal-left", "", "A comma separated list values to inverse match the left branch on. eg: paycode:salary")
var inverseLiteralRightMatchInput = flag.String("inverse-match-literal-right", "", "A comma separated list values to inverse match the right branch on. eg: paycode:extra_hours")
var groupKey = flag.String("group-key", "id", "The header with which to group (sorted!) rows")

var logger = util.Logger{}

func main() {
	flag.Parse()

	if *version {
		logger.Info(currentVersion)
		os.Exit(0)
	}

	output := csv.NewWriter(os.Stdout)

	matchOn := []matchSet{}

	if *matchInput != "" {
		columns := strings.Split(*matchInput, ",")
		for _, set := range columns {
			bits := strings.SplitN(set, ":", 2)

			matchOn = append(matchOn, matchSet{
				Left:  bits[0],
				Right: bits[1],
			})
		}
	}

	if *inverseMatchInput != "" {
		columns := strings.Split(*inverseMatchInput, ",")
		for _, set := range columns {
			bits := strings.SplitN(set, ":", 2)

			matchOn = append(matchOn, matchSet{
				Inverse: true,
				Left:    bits[0],
				Right:   bits[1],
			})
		}
	}

	if *literalLeftMatchInput != "" {
		columns := strings.Split(*literalLeftMatchInput, ",")
		for _, set := range columns {
			bits := strings.SplitN(set, ":", 2)

			matchOn = append(matchOn, matchSet{
				LiteralLeft: true,
				Left:        bits[0],
				Right:       bits[1],
			})
		}
	}

	if *literalRightMatchInput != "" {
		columns := strings.Split(*literalRightMatchInput, ",")
		for _, set := range columns {
			bits := strings.SplitN(set, ":", 2)

			matchOn = append(matchOn, matchSet{
				LiteralRight: true,
				Left:         bits[0],
				Right:        bits[1],
			})
		}
	}

	if *inverseLiteralLeftMatchInput != "" {
		columns := strings.Split(*inverseLiteralLeftMatchInput, ",")
		for _, set := range columns {
			bits := strings.SplitN(set, ":", 2)

			matchOn = append(matchOn, matchSet{
				Inverse:     true,
				LiteralLeft: true,
				Left:        bits[0],
				Right:       bits[1],
			})
		}
	}

	if *inverseLiteralRightMatchInput != "" {
		columns := strings.Split(*inverseLiteralRightMatchInput, ",")
		for _, set := range columns {
			bits := strings.SplitN(set, ":", 2)

			matchOn = append(matchOn, matchSet{
				Inverse:      true,
				LiteralRight: true,
				Left:         bits[0],
				Right:        bits[1],
			})
		}
	}

	if err := ducky(os.Stdin, *output, matchOn, *groupKey); err != nil {
		logger.Fatal(err)
	}

	output.Flush()

	logDone()
}

func ducky(input io.Reader, output csv.Writer, matchOn []matchSet, groupKey string) error {
	work, errors := util.ReadSourceAsync(input)

	var cachedErr error
	go (func() {
		for err := range errors {
			logger.Error(err)
			cachedErr = err
		}
	})()

	headersPrinted := false

	prevLine := util.Line{}
	group := []util.Line{}

	for line := range work {
		if !headersPrinted {
			if err := output.Write(append(line.Headers, "ducky_taped")); err != nil {
				return err
			}
			headersPrinted = true
		}

		if prevLine.Data[groupKey] == line.Data[groupKey] || len(prevLine.Headers) == 0 {
			group = append(group, line)
		} else {
			result := matchGroup(group, matchOn)
			for _, l := range result {
				if err := emitLine(l, &output); err != nil {
					return err
				}
			}
			group = []util.Line{line}
		}
		prevLine = line
	}

	result := matchGroup(group, matchOn)
	for _, l := range result {
		if err := emitLine(l, &output); err != nil {
			return err
		}
	}

	output.Flush()

	return cachedErr
}

func doLinesMatch(left util.Line, right util.Line, match matchSet) bool {
	matched := false

	if match.LiteralRight {
		if right.Data[match.Left] == match.Right {
			matched = true
		}
	} else if match.LiteralLeft {
		if left.Data[match.Left] == match.Right {
			matched = true
		}
	} else {
		if left.Data[match.Left] == right.Data[match.Right] {
			matched = true
		}
	}

	if match.Inverse {
		matched = !matched
	}

	return matched
}

func anyLinesMatch(group []util.Line, match matchSet) bool {
	anyMatch := false
	group = append(group, util.Line{})
	for i, line := range group {
		left := util.Line{}
		if i != 0 {
			left = group[i-1]
		}
		matched := doLinesMatch(left, line, match)
		if matched {
			//if doLinesMatch(group[i-1], line, match) {
			anyMatch = true
		}
	}
	return anyMatch
}

func matchGroup(group []util.Line, allMatches []matchSet) []util.Line {
	if len(group) == 1 {
		group[0].Data["ducky_taped"] = "false"
		return group
	}

	for _, match := range allMatches {
		if match.MatchAny {
			if !anyLinesMatch(group, match) {
				for _, line := range group {
					line.Data["ducky_taped"] = "false"
				}
				return group
			}
		}
	}

	matchOn := []matchSet{}
	for _, match := range allMatches {
		if !match.MatchAny {
			matchOn = append(matchOn, match)
		}
	}
	requiredMatches := len(matchOn)

	prevLine := util.Line{}
	result := []util.Line{}

	for i, line := range group {
		numMatches := 0
		if i == 0 {
			prevLine = line
			prevLine.Data["ducky_taped"] = "false"
			continue
		}
		for _, match := range matchOn {
			matched := doLinesMatch(prevLine, line, match)
			if matched {
				numMatches += 1
			}
		}

		// If we hit all the matches for this line, merge it into prevLine and discard it
		if numMatches == requiredMatches {
			for _, match := range matchOn {
				prevLine.Data["ducky_taped"] = "true"
				prevLine.Data[match.Left] = line.Data[match.Left]
			}
		} else {
			result = append(result, prevLine)
			prevLine = line
			prevLine.Data["ducky_taped"] = "false"
		}
	}

	return append(result, prevLine)
}

func emitLine(line util.Line, output *csv.Writer) error {
	newLine := []string{}
	line.Headers = append(line.Headers, "ducky_taped")
	for _, header := range line.Headers {
		newLine = append(newLine, line.Data[header])
	}

	return output.Write(newLine)
}

type matchSet struct {
	Left         string
	LiteralLeft  bool
	Right        string
	LiteralRight bool
	Inverse      bool
	MatchAny     bool
}

func logDone() {
	if *quiet {
		return
	}
	logger.Info(`
  _      _      _
>(.)__ <(.)__ =(.)__
 (___/  (___/  (___/
`)
}
