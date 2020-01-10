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

var logger = util.Logger{}

func main() {
	flag.Parse()

	if *version {
		logger.Info(currentVersion)
		os.Exit(0)
	}

	if err := flimflam(os.Stdin, os.Stdout); err != nil {
		logger.Error(err)
	}

	logDone()
}

func flimflam(input io.Reader, output io.Writer) error {

	reader := csv.NewReader(input)
	line, err := reader.Read()
	if err != nil {
		return err
	}

	output.Write([]byte(strings.Join(line, ":STRING,")))
	output.Write([]byte(":STRING\n"))

	return nil
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
