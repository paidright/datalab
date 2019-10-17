package util

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type lineFunc func(line map[string]string, headers []string, lineNumber int) error

type Line struct {
	Data    map[string]string
	Headers []string
	Number  int
}

func ReadSourceAsync(input io.Reader) (chan Line, chan error) {
	r := csv.NewReader(input)

	cols := []string{}
	lineNumber := 0

	work := make(chan Line, 1000)
	errors := make(chan error)

	go (func() {
		for {
			lineNumber += 1
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errors <- err
				break
			}

			if lineNumber == 1 {
				cols = record
				continue
			}

			line := Line{
				Number:  lineNumber,
				Headers: cols,
				Data:    map[string]string{},
			}

			for i, col := range cols {
				if !(len(record) > i) {
					errors <- fmt.Errorf("WARN Invalid CSV. Missing column at line %d", i)
					break
				}
				line.Data[col] = record[i]
			}

			work <- line
		}

		close(work)
		close(errors)
	})()

	return work, errors
}

func ReadFileAsync(path string) (chan Line, chan error) {
	f, err := os.Open(path)

	if err != nil {
		work := make(chan Line)
		errors := make(chan error)
		errors <- err
		return work, errors
	}

	work, errors := ReadSourceAsync(f)

	return work, errors
}

func ReadSource(input io.Reader, handler lineFunc) error {
	r := csv.NewReader(input)

	cols := []string{}
	lineNumber := 0

	for {
		lineNumber += 1
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if lineNumber == 1 {
			cols = record
			continue
		}

		line := map[string]string{}
		for i, col := range cols {
			if !(len(record) > i) {
				return fmt.Errorf("WARN Invalid CSV. Missing column at line %d", i)
			}
			line[col] = record[i]
		}

		if err := handler(line, cols, lineNumber); err != nil {
			return err
		}
	}

	return nil
}

func ReadFile(path string, handler lineFunc) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	err = ReadSource(f, handler)
	if err != nil {

		if err := f.Close(); err != nil {
			log.Println("ERROR closing file: ", err)
		}
	}

	if err != nil {
		return fmt.Errorf("error in file %s: %w", path, err)
	}

	return err
}

func ReadHeadersFromSource(input io.Reader) ([]string, error) {
	r := csv.NewReader(input)

	return r.Read()
}

func ReadHeaders(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return []string{}, err
	}

	return ReadHeadersFromSource(f)
}

func ListFiles(inputPath string, ignored []string, targetExtensions []string) ([]string, error) {
	files, err := ioutil.ReadDir(inputPath)
	if err != nil {
		return []string{}, err
	}

	names := []string{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if Contains(file.Name(), ignored) {
			continue
		}
		ext := path.Ext(file.Name())
		if Contains(ext, targetExtensions) {
			names = append(names, file.Name())
		}
	}

	return names, nil
}

func Contains(target string, considered []string) bool {
	for _, item := range considered {
		if item == target {
			return true
		}
	}

	return false
}
