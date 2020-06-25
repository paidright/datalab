package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/paidright/datalab/util"
)

var version = flag.Bool("version", false, "Just print the version and exit")
var quiet = flag.Bool("quiet", false, "Tone down the output noise")
var format = flag.String("format", "kv", "Output in key-value or json format. Valid: kv, json")

var logger = util.Logger{}

func main() {
	flag.Parse()

	if *version {
		logger.Info(currentVersion)
		os.Exit(0)
	}

	if err := flimflam(os.Stdin, *format, os.Stdout); err != nil {
		logger.Error(err)
	}

	logDone()
}

func flimflam(input io.Reader, format string, output io.Writer) error {
	reader := csv.NewReader(input)
	line, err := reader.Read()
	if err != nil {
		return err
	}

	switch format {
	case "kv":
		if _, err := output.Write([]byte(strings.Join(line, ":STRING,"))); err != nil {
			return err
		}
		if _, err := output.Write([]byte(":STRING\n")); err != nil {
			return err
		}
	case "json":
		schema := []coldef{}

		for _, col := range line {
			schema = append(schema, coldef{
				Name: col,
				Type: "STRING",
			})
		}

		result, err := json.Marshal(schema)
		if err != nil {
			return err
		}
		if _, err := output.Write(result); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Invalid format specified")
	}

	return nil
}

type coldef struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func logDone() {
	if *quiet {
		return
	}
	logger.Info(`
      .~~~~'\~~\
     ;       ~~ \
     |           ;
 ,--------,______|---.
/          \-----'    \
'.__________'-_______-'
`)
}
