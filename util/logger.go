package util

import (
	"log"
	"os"
	"strings"
)

type Logger struct{}

func (_ Logger) Info(input ...string) {
	log.Println(strings.Join(append([]string{"INFO", os.Args[0]}, input...), " "))
}

func (_ Logger) Fatal(err error) {
	log.Fatal(strings.Join(append([]string{"ERROR", os.Args[0]}, err.Error()), " "))
}

func (_ Logger) Error(err error) {
	log.Println(strings.Join(append([]string{"ERROR", os.Args[0]}, err.Error()), " "))
}
