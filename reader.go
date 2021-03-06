// This file contains the reader, reading an OSM-file (usually .osm or .xml
// files) and send each changeset as one-line string to a given channel.
package main

import (
	"bufio"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hauke96/osm-changeset-crawler/common"
	"github.com/hauke96/sigolo"
)

// read the given file and cache "cacheSize" many changesets from that file
// before handing it over to the "changesetStringChan". The pipeline receives
// an array of strings, each string is one changeset.
func read(fileName string, changesetStringChan chan<- []string, finishWaitGroup *sync.WaitGroup) {
	defer close(changesetStringChan)
	defer finishWaitGroup.Done()

	clock := time.Now()

	changesetPrefix := "<changeset "
	changesetSuffix := "</changeset>"
	changesetOneLineSuffix := "/>"

	cache := make([]string, common.CACHE_SIZE)

	readChangesetSets := 0
	var line string

	// Open file
	fileHandle, err := os.Open(fileName)
	sigolo.FatalCheck(err)
	defer fileHandle.Close()
	sigolo.Info("Opened file")

	const capacity = 64 * 1024 * 1024
	buf := make([]byte, capacity)
	scanner := bufio.NewScanner(fileHandle)
	scanner.Buffer(buf, capacity)

	sigolo.Info("Created scanner")

	for scanner.Scan() {
		clock = time.Now()

		for i := 0; i < common.CACHE_SIZE && scanner.Scan(); i++ {
			line = strings.TrimSpace(scanner.Text())

			// New changeset starts
			if strings.HasPrefix(line, changesetPrefix) {
				// Read all lines of this changeset
				changesetString := line

				// If the read line is not a one-line-changeset like
				// "<changeset id=123 open=false ... />"), then read the other lines
				if !strings.HasSuffix(changesetString, changesetOneLineSuffix) {
					for scanner.Scan() {
						line = strings.TrimSpace(scanner.Text())
						changesetString += line

						// Changeset ends
						if strings.HasPrefix(line, changesetSuffix) {
							break
						}
					}
				}

				// Done reading the changeset, add it to the cache
				cache[i] = changesetString
			}
		}

		sigolo.Info("Read changeset set %d", readChangesetSets)
		sigolo.Info("Reading took %dms", time.Since(clock).Milliseconds())

		changesetStringChan <- cache
		cache = make([]string, common.CACHE_SIZE)

		sigolo.Info("Total reoundtrip time was %dms", time.Since(clock).Milliseconds())
	}

	sigolo.Debug("Reading finished, send remaining strings")

	changesetStringChan <- cache
}
