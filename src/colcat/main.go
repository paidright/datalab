package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/paidright/datalab/util"
)

var version = flag.Bool("version", false, "Just print the version and exit")
var quiet = flag.Bool("quiet", false, "Tone down the output noise")
var targetFile = flag.String("target_file", "targets.csv", "The file containing the list of targets")

func main() {
	flag.Parse()
	if *version {
		log.Println(currentVersion)
		os.Exit(0)
	}

	targets, err := readTargets(*targetFile)
	if err != nil {
		log.Fatal(err)
	}

	output := csv.NewWriter(os.Stdout)

	if err := processFile(os.Stdin, targets, *output); err != nil {
		log.Fatal("ERROR", err)
	}

	output.Flush()

	logDone()
}

func processFile(input io.Reader, targets []target, output csv.Writer) error {
	cachedPassThroughHeaders := []string{}

	handleHeaders := func(headers []string) ([]string, error) {
		if len(cachedPassThroughHeaders) > 0 {
			return cachedPassThroughHeaders, nil
		}

		if err := validateTargets(targets, headers); err != nil {
			return []string{}, err
		}

		allSources := []string{}
		for _, target := range targets {
			allSources = append(allSources, target.sources...)
		}

		passthroughHeaders := []string{}
		for _, header := range headers {
			if !util.Contains(header, allSources) {
				passthroughHeaders = append(passthroughHeaders, header)
			}
		}

		allDests := []string{}
		for _, target := range targets {
			allDests = append(allDests, target.dest)
		}

		cachedPassThroughHeaders = passthroughHeaders

		newHeaders := append(passthroughHeaders, allDests...)

		if err := output.Write(newHeaders); err != nil {
			return []string{}, err
		}
		output.Flush()

		return passthroughHeaders, nil
	}

	work, errors := util.ReadSourceAsync(input)

	for line := range work {
		passthroughHeaders, err := handleHeaders(line.Headers)
		if err != nil {
			return err
		}

		result := []string{}
		for _, header := range passthroughHeaders {
			result = append(result, line.Data[header])
		}
		for _, target := range targets {
			destCell := []string{}
			for _, source := range target.sources {
				destCell = append(destCell, line.Data[source])
			}
			result = append(result, strings.Join(destCell, target.sep))
		}

		err = output.Write(result)
		output.Flush()
		if err != nil {
			return err
		}
	}

	var cachedErr error

	for err := range errors {
		log.Println("ERROR", err)
		cachedErr = err
	}

	return cachedErr
}

func validateTargets(targets []target, headers []string) error {
	for _, target := range targets {
		for _, source := range target.sources {
			if !util.Contains(source, headers) {
				return fmt.Errorf("target header %s does not exist in input CSV", target)
			}
		}
	}

	return nil
}

func readTargets(filePath string) ([]target, error) {
	targets := []target{}

	// Read in the target input file
	if err := util.ReadFile(filePath, func(line map[string]string, headers []string, lineNumber int) error {
		targets = append(targets, target{
			sources: strings.Split(line["sources"], ":"),
			sep:     line["sep"],
			dest:    line["dest"],
		})
		return nil
	}); err != nil {
		return targets, err
	}

	return targets, nil
}

type target struct {
	sources []string
	dest    string
	sep     string
}

func logDone() {
	if *quiet {
		return
	}
	log.Println(`
      /\_____/\
     /  o   o  \
    ( ==  ^  == )
     )_________(
    (  |Colin|  )
   ( (  )   (  ) )
  (__(__)___(__)__)
  G'day. Name's Col.`)
}
