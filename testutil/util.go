package testutil

import (
	"io/ioutil"
	"log"
)

func Load(path string) []byte {
	filePath := "./testdata/" + path
	buffer, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("cannot load data: %s", filePath)
	}
	return buffer
}
