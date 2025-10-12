package main

import (
	"flag"
	"fmt"
)

func main() {
	n := flag.Int("n", 10, "number of lines")
	// follow := flag.Bool("f", false, "follow the file as it grows")

	flag.Parse()
	files := flag.Args()

	for _, filePath := range files {
		fr, err := NewFileReader(filePath)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", filePath, err)
			continue
		}
		iterator, err := fr.Tail(*n)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filePath, err)
			continue
		}
		defer iterator.Close()

		value, err := iterator.Next()
		if err != nil {
			fmt.Printf("Error reading from iterator: %v\n", err)
			continue
		}
		for len(value) > 0 {
			fmt.Printf("%s", value)
			value, err = iterator.Next()
		}
	}
}
