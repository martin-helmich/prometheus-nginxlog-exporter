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
	"runtime/pprof"
	"sync"
)

// SetupCPUProfiling starts CPU profiling if an outputFile is specified
func SetupCPUProfiling(outputFile string, stopChan <-chan bool, stopHandlers *sync.WaitGroup) {
	if outputFile == "" {
		return
	}

	f, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("writing CPU profile to file %s\n", outputFile)

	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}

	stopHandlers.Add(1)

	go func() {
		<-stopChan

		fmt.Printf("stopping CPU profiling...\n")
		pprof.StopCPUProfile()

		stopHandlers.Done()
	}()
}
