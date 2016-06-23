package gonx

import (
	"bufio"
	"bytes"
	"io"
	"sync"
)

func handleError(err error) {
	//fmt.Fprintln(os.Stderr, err)
}

// Iterate over given file and map each it's line into Entry record using
// parser and apply reducer to the Entries channel. Execution terminates
// when result will be readed from reducer's output channel, but the mapper
// works and fills input Entries channel until all lines will be read from
// the fiven file.
func MapReduce(file io.Reader, parser StringParser, reducer Reducer) chan *Entry {
	// Input file lines. This channel is unbuffered to publish
	// next line to handle only when previous is taken by mapper.
	var lines = make(chan string)

	// Host thread to spawn new mappers
	var entries = make(chan *Entry, 10)
	go func(topLoad int) {
		// Create semafore channel with capacity equal to the output channel
		// capacity. Use it to control mapper goroutines spawn.
		var sem = make(chan bool, topLoad)
		for i := 0; i < topLoad; i++ {
			// Ready to go!
			sem <- true
		}

		var wg sync.WaitGroup
		for {
			// Wait until semaphore becomes available and run a mapper
			if !<-sem {
				// Stop the host loop if false received from semaphore
				break
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				// Take next file line to map. Check is channel closed.
				line, ok := <-lines
				// Return immediately if lines channel is closed
				if !ok {
					// Send false to semaphore channel to indicate that job's done
					sem <- false
					return
				}
				entry, err := parser.ParseString(line)
				if err == nil {
					// Write result Entry to the output channel. This will
					// block goroutine runtime until channel is free to
					// accept new item.
					entries <- entry
				} else {
					handleError(err)
				}
				// Increment semaphore to allow new mapper workers to spawn
				sem <- true
			}()
		}
		// Wait for all mappers to complete, then send a quit signal
		wg.Wait()
		close(entries)
	}(cap(entries))

	// Run reducer routine.
	var output = make(chan *Entry)
	go reducer.Reduce(entries, output)

	go func() {
		reader := bufio.NewReader(file)
		line, err := readLine(reader)
		for err == nil {
			// Read next line from the file and feed mapper routines.
			lines <- line
			line, err = readLine(reader)
		}
		close(lines)

		if err != nil && err != io.EOF {
			handleError(err)
		}
	}()

	return output
}

func readLine(reader *bufio.Reader) (string, error) {
	line, isPrefix, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	if !isPrefix {
		return string(line), nil
	}
	var buffer bytes.Buffer
	_, err = buffer.Write(line)
	for isPrefix && err == nil {
		line, isPrefix, err = reader.ReadLine()
		if err == nil {
			_, err = buffer.Write(line)
		}
	}
	return buffer.String(), err
}
