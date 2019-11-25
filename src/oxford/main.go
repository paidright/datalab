package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
)

var delim = flag.String("delimiter", "", "The delimiter currently used by the input data")

func main() {
	flag.Parse()
	output := csv.NewWriter(os.Stdout)

	if err := oxford(os.Stdin, output, []rune(*delim)[0]); err != nil {
		log.Fatal(err)
	}

	output.Flush()
}

func oxford(input io.Reader, output *csv.Writer, delim rune) error {
	r := csv.NewReader(input)

	r.Comma = delim

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := output.Write(record); err != nil {
			return err
		}
	}
	return nil
}
