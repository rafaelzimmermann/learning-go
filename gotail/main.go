package main

import (
	"flag"
	"fmt"
)

func main() {
	n := flag.Int("n", 10, "number of lines")
	follow := flag.Bool("f", false, "follow the file as it grows")

	flag.Parse()
	files := flag.Args()

	fmt.Printf("n=%d, f=%t, files=%q\n", *n, *follow, files)

	for _, filePath := range files {
		fr, err := NewFileReader(filePath)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", filePath, err)
			continue
		}
		lines, err := fr.Tail(*n)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filePath, err)
			continue
		}
		for _, line := range lines {
			fmt.Println(line)
		}
	}
}
