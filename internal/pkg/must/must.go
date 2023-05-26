package must

import (
	"log"
	"os"
)

func Open(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	return file

}

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
