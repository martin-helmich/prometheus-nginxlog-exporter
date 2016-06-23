// Example program that reads big nginx file from stdin line by line
// and measure reading time. The file should be big enough, at least 500K lines
package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var count int
	start := time.Now()
	for scanner.Scan() {
		// A dummy action, jest read line by line
		scanner.Text()
		count++
	}
	duration := time.Since(start)
	fmt.Printf("%v lines readed, it takes %v\n", count, duration)
}
