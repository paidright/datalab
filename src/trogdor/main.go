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
var columns = flag.String("columns", "User_ID,Effective_Start_Date", "The column(s) to shuffle to the start of the file")

func main() {
	flag.Parse()
	if *version {
		log.Println(currentVersion)
		os.Exit(0)
	}

	output := csv.NewWriter(os.Stdout)

	if err := trogdor(os.Stdin, *columns, output); err != nil {
		log.Fatal("ERROR", err)
	}

	output.Flush()

	logDone()
}

func trogdor(input io.Reader, cols string, output *csv.Writer) error {
	cachedNewHeaders := []string{}
	targets := []string{}

	processHeaders := func(headers []string) ([]string, error) {
		if len(cachedNewHeaders) > 0 {
			return cachedNewHeaders, nil
		}

		innerTargets, err := validateTargets(cols, headers)
		if err != nil {
			return []string{}, err
		}

		targets = innerTargets

		newHeaders := []string{}
		for _, header := range headers {
			if !contains(header, targets) {
				newHeaders = append(newHeaders, header)
			}
		}

		output.Write(append(targets, newHeaders...))
		output.Flush()

		cachedNewHeaders = newHeaders

		return newHeaders, nil
	}

	work, errors := util.ReadSourceAsync(input)

	for line := range work {
		if _, err := processHeaders(line.Headers); err != nil {
			return err
		}

		result := []string{}
		for _, target := range targets {
			result = append(result, line.Data[target])
		}

		for _, header := range line.Headers {
			if contains(header, targets) {
				continue
			}
			result = append(result, line.Data[header])
		}

		if err := output.Write(result); err != nil {
			return err
		}
		output.Flush()
	}

	var cachedErr error

	for err := range errors {
		log.Println("ERROR", err)
		cachedErr = err
	}

	return cachedErr
}

func validateTargets(cols string, headers []string) ([]string, error) {
	targets := strings.Split(cols, ",")

	for _, target := range targets {
		found := false
		for _, header := range headers {
			if found {
				continue
			}
			if header == target {
				found = true
			}
		}
		if !found {
			return []string{}, fmt.Errorf("target header %s does not exist in input CSV", target)
		}
	}

	return targets, nil
}

func contains(term string, pool []string) bool {
	for _, v := range pool {
		if v == term {
			return true
		}
	}
	return false
}

func logDone() {
	if *quiet {
		return
	}
	log.Println(`
and the trogdor comes in the NIIGGHHHTTTTT!!!!!!!!
                                                 :::
                                             :: :::.
                       \/,                    .:::::
           \),          \'-._                 :::888
           /\            \   '-.             ::88888
          /  \            | .(                ::88
         /,.  \           ; ( '              .:8888
            ), \         / ;''               :::888
           /_   \     __/_(_                  :88
             '. ,'..-'      '-._    \  /      :8
               )__ '.           '._ .\/.
              /   '. '             '-._______m         _,
  ,-=====-.-;'                 ,  ___________/ _,-_,'"'/__,-.
 C   =--   ;                   '.'._    V V V       -=-'"#==-._
:,  \     ,|      UuUu _,......__   '-.__A_A_ -. ._ ,--._ ",'' '-
||  |'---' :    uUuUu,'          ''--...____/   '" '".   '
|'  :       \   UuUu:
:  /         \   UuUu'-._
 \(_          '._  uUuUu '-.
 (_3             '._  uUu   '._
                    ''-._      '.
                         '-._    '.
                             '.    \
                               )   ;
                              /   /
               '.        |\ ,'   /
                 ",_A_/\-| '   ,'
                   '--..,_|_,-'\
                          |     \
                          |      \__
                          |__`)
}
