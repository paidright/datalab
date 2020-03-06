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
var literalLeftMatchInput = flag.String("match-left-literal", "", "A comma separated list values to match the left branch on. eg: paycode:salary")
var literalRightMatchInput = flag.String("match-right-literal", "", "A comma separated list values to match the right branch on. eg: paycode:extra_hours")

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
			bits := strings.Split(set, ":")

			matchOn = append(matchOn, matchSet{
				Left:  bits[0],
				Right: bits[1],
			})
		}
	}

	if *literalLeftMatchInput != "" {
		columns := strings.Split(*literalLeftMatchInput, ",")
		for _, set := range columns {
			bits := strings.Split(set, ":")

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
			bits := strings.Split(set, ":")

			matchOn = append(matchOn, matchSet{
				LiteralRight: true,
				Left:         bits[0],
				Right:        bits[1],
			})
		}
	}

	if err := ducky(os.Stdin, *output, matchOn); err != nil {
		logger.Fatal(err)
	}

	output.Flush()

	logDone()
}

func ducky(input io.Reader, output csv.Writer, matchOn []matchSet) error {
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

	for line := range work {
		line.Headers = append(line.Headers, "ducky_taped")
		line.Data["ducky_taped"] = "false"

		if !headersPrinted {
			if err := output.Write(line.Headers); err != nil {
				return err
			}
			headersPrinted = true
		}

		if len(prevLine.Headers) == 0 {
			prevLine = line
			continue
		}

		numMatches := 0
		for _, match := range matchOn {
			if match.LiteralRight {
				if line.Data[match.Left] == match.Right {
					numMatches += 1
				}
			} else if match.LiteralLeft {
				if prevLine.Data[match.Left] == match.Right {
					numMatches += 1
				}
			} else {
				if prevLine.Data[match.Left] == line.Data[match.Right] {
					numMatches += 1
				}
			}
		}

		if numMatches == len(matchOn) {
			for _, match := range matchOn {
				prevLine.Data["ducky_taped"] = "true"
				prevLine.Data[match.Left] = line.Data[match.Left]
			}
		} else {
			if err := emitLine(prevLine, &output); err != nil {
				return err
			}

			prevLine.Data = line.Data
		}
	}

	if err := emitLine(prevLine, &output); err != nil {
		return err
	}

	return cachedErr
}

func emitLine(line util.Line, output *csv.Writer) error {
	newLine := []string{}
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
