// +build munge_cwd

package main

import (
	"log"
	"os"
	"path/filepath"
)

func init() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("INFO Darwin platform. Changing CWD to: ", ex)

	exPath := filepath.Dir(ex)
	os.Chdir(exPath)
}
