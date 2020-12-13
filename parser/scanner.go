package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Scan for input from stdin
func Scan(output chan<- string) {
	defer close(output)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		output <- strings.TrimSpace(line)
	}

	err := scanner.Err()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	} else {
		fmt.Fprint(os.Stdout, "No more input")
	}
}
