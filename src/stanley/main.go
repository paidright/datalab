package main

import (
	"datalab/util"
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"strconv"
)

var version = flag.Bool("version", false, "Just print the version and exit")
var quiet = flag.Bool("quiet", false, "Tone down the output noise")
var left = flag.String("left", "left.csv", "The file containing the left hand side of the join")
var joinkey = flag.String("join-key", "id", "The column on which to do the join")

func main() {
	flag.Parse()
	if *version {
		log.Println(currentVersion)
		os.Exit(0)
	}

	log.Printf("INFO Stanley is inner joining %s with stdin on %s\n", *left, *joinkey)

	output := csv.NewWriter(os.Stdout)

	leftInput, err := os.Open(*left)
	if err != nil {
		log.Fatal(err)
	}

	if err := join(*joinkey, leftInput, os.Stdin, output); err != nil {
		log.Fatal(err)
	}

	output.Flush()

	logDone()
}

func join(key string, left io.Reader, right io.Reader, dest *csv.Writer) error {
	cachedHeaders := []string{}
	leftHeaders := []string{}

	handleHeaders := func(cols []string) ([]string, error) {
		if len(cachedHeaders) > 0 {
			return cachedHeaders, nil
		}

		headers := util.Uniq(append(leftHeaders, cols...))

		if err := dest.Write(headers); err != nil {
			return []string{}, err
		}
		dest.Flush()

		cachedHeaders = headers

		return cachedHeaders, nil
	}

	leftCache := map[string]map[string]string{}

	err := util.ReadSource(left, func(line map[string]string, cols []string, lineNumber int) error {
		leftHeaders = cols
		line["left_original_line_number"] = strconv.Itoa(lineNumber)
		leftCache[line[key]] = line
		return nil
	})

	if err != nil {
		return err
	}

	work, errors := util.ReadSourceAsync(right)

	for line := range work {
		headers, err := handleHeaders(line.Headers)
		if err != nil {
			return err
		}

		for k, v := range leftCache[line.Data[key]] {
			line.Data[k] = v
		}
		line.Data["right_original_line_number"] = strconv.Itoa(line.Number)

		output := []string{}
		for _, header := range headers {
			output = append(output, line.Data[header])
		}

		if err = dest.Write(output); err != nil {
			return err
		}
		dest.Flush()
	}

	var cachedErr error

	for err := range errors {
		log.Println("ERROR", err)
		cachedErr = err
	}

	return cachedErr
}

func logDone() {
	if *quiet {
		return
	}
	log.Println(`
                                                                                                        @@@@@@@@@
                                                                                                     @@@@@@@@@@@@@@
                                                                                                   @@@@@@@@@@@@@@@     @@@@@@@@
                                                                                                 @@@@@@@@@@@@@@@@ @@@@@@@@@@@@@
                                                                                               @@@@@@@@  @@@@@@@@@@@@@@ @@@@@@
                                                                                             @@@@@@@@  @@@@@@@@@@@@  @@@@@@
                                                                                        @@@@@@@@@@@@  @@@@@@@@@   @@@@@@@@
                                                                                      @@@@@@@@@@@@@ @@@@@@@@@@ @@@@@@@@@@
                                                                                    @@@@@@@@@@@@@ @@@@@@@@@@@@@@@@@@@@@@
                                                                                   @@@@@@@@@@@@ @@@@@@@@@@  @@@@@@@@@@@@
                                                                                   @@@@@@@@@@@ @@@@@@@@@@    @@@@@@@@@@@
                                                                                   @@@@@@@@@ @@@@@@@@@@@      @@@@@@@@@ @
                                                                               @@@@@@@@@@@@@@@@@@@@@@@ @@      @@@@@@@@@@@@
                                                                             @@@@@@@@@@@@@@@@@@@@@@@@   @@      @@@@@@@@@@@@
                                                                            @@@@@@@@@@@@@@@@@@@@@@@       @      @@@@@@@@@@@@@
                                                                          @@@@@@@@@@@@@@@@@@@@@@@@         @      @@@@@@@@@@@@@
                                                                         @@@@@@@@@@@@ @@@@@@@@@@            @      @@@@@@@@@@@ @
                                                                        @@@@@@@@@@@@@@@@@@@@@@@              @       @@@@@@@@@@@@
                                                                      @@@@@@@@@@@@@@@@@@@@@@@                 @@      @@@@@@@@@@@@
                                                                    @@@@@@@@@@@@@@@@@@@@@@@@@                   @      @@@@@@@@@@@@
                      @@@ @ @@@@@@                                @@@@@@@@@@@@@@@@@@@@@@@@@                      @       @@@@@@@@@@@
                   @@@@@@@@@@@@@@@@@@@                          @@@@@@@@@@@@@@@@@@@@@@@@@@  @@                    @       @@@@@@@@@@
                 @@@@@@@@@@@@@@@@@@@@@@                      @@@@@@@@@@@@@@@@@@@@@@@@@@@@@   @@                    @      @@@@@@@@@@@
                 @@@@@@@@@@ @@@@@@@  @@@                @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ @@@  @@@                   @@      @@@@@@@@@@
                 @@@@@@@@@@@ @@@@@@@@@@@           @@@@  @@@@@@@   @@@@@@@@@@@@@@@@@@@@@@    @@@                   @    @ @@@@@@@@@@@@
                 @@@@ @@@@@@   @@@@@@@@@       @@@@@@@  @@@@@@ @@@@@@@@@@@@@@@@@@@@@@@ @@@  @@@@                   @@@ @@ @@@@@@@@@@@
                @@@@@@@@@@@      @@@@@@      @@@@ @@    @@@@@@@@@@@@@@@@ @@@@@@@@@@@@ @    @@  @                 @@@@@ @ @@@@@@@@@@@@@
                 @@@@@@@@@@       @@@@@    @@          @ @@@@@  @@@@@@@@@@@@@@@@@@@@@@@@@  @                 @@@@@@   @@@@@@@@@@@@@@@@@@@@ @
                 @@@@@  @                @@      @    @@@@@@@  @ @@@@@@@@@@@@@@@@@@@@@@@@@@@@@     @@@@@ @@ @@@@     @@@@@@@@@@@@@@@@@@@@@@ @@@@  @
                  @@@                 @@@             @@ @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@    @   @@@@@ @@@@@                  @@@@@@@@@@@@@@@@@@@@@@@@  
                    @               @@               @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@              @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
                     @             @   @@@          @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
                     @@@@@@            @  @          @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
                     @@@@@@   @@@@           @@@@@  @ @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
                    @ @@@@@   @@@@               @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
       @@@@@        @@@@@@@@ @@@@ @  @@@@      @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
    @@              @@@@@@@ @@@@@    @    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
   @       @@@@@@@@     @@          @ @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
  @@   @@@  @@@ @@ @@@@@@@@@ @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@  @ @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@    @@@@@@@@@   @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@     @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
   @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@`)
}
