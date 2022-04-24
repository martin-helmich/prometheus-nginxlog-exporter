/*
 * Copyright 2019 Martin Helmich <martin@helmich.me>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package prof

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
)

// SetupMemoryProfiling starts memory profiling if an outputFile is specified
func SetupMemoryProfiling(outputFile string, stopChan <-chan bool, stopHandlers *sync.WaitGroup) {
	if outputFile == "" {
		return
	}

	runtime.MemProfileRate = 1
	stopHandlers.Add(1)

	go func() {
		<-stopChan

		f, err := os.Create(outputFile)
		if err != nil {
			panic(err)
		}

		fmt.Printf("writing memory profile to file %s\n", outputFile)

		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			panic(err)
		}

		f.Close()
		stopHandlers.Done()
	}()
}
