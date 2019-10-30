// Package regexmachine helps with processing regex rules of streams of data and looking for rule matches
package regexmachine

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"regexp"
)

// RuleSet Rules to check against
type RuleSet []*regexp.Regexp

// Public regexes
var (
	Email = regexp.MustCompile(`[A-Za-z0-9_.]+((\ ?(\[|\()?\ ?@\ ?(\)|\])?\ ?)|(\ ?(\[|\()\ ?[aA][tT]\ ?(\)|\])\ ?))[0-9a-z.-]+`)

	// DefaultSet
	DefaultRules = RuleSet{Email}
)

// -- Functions on RuleSet --

// GetMatchedRules Given bytes return all regexes that match
func (rules *RuleSet) GetMatchedRules(data *[]byte) []*regexp.Regexp {
	matched := []*regexp.Regexp{}
	for _, rule := range *rules {
		if rule.Match(*data) {
			matched = append(matched, rule)
		}
	}

	return matched
}

// MatchesRule Given bytes return if any rule matches
func (rules *RuleSet) MatchesRule(data *json.RawMessage) bool {
	for _, rule := range *rules {
		if rule.Match(*data) {
			return true
		}
	}

	return false
}

// MatchesRuleReader Given a reader, return true if any rule matches in the stream.  Will read ENTIRE READER
// Spawns multiple go routines to check each rule
// Use limit reader to prevent reading forever
func (rules *RuleSet) MatchesRuleReader(ctx context.Context, reader io.Reader) bool {
	foundMatch := make(chan bool)
	finishedWorkers := make(chan bool)
	rulesLen := len(*rules)

	workerContext, cancelWorkers := context.WithCancel(ctx)

	// We need to duplicate the reader stream for each worker
	// Create reader and writer for each worker thread
	workerWriters := []*io.PipeWriter{}
	workerReaders := []*io.PipeReader{}
	for range *rules {
		r, w := io.Pipe()
		workerWriters = append(workerWriters, w)
		workerReaders = append(workerReaders, r)
	}

	// Read stream duplicator to write to all workerWriters
	go func() {
		// Defer closing all workers
		closeAll := func() {
			for _, writer := range workerWriters {
				writer.Close()
			}
		}
		defer closeAll()

		for {
			// Read
			buf := make([]byte, 1024)
			if _, err := reader.Read(buf); err != nil {
				return
			}

			// Writer to each worker stream
			for _, writer := range workerWriters {
				writer.Write(buf)
			}

			// Check if we are canceled
			select {
			case <-workerContext.Done():
				return
			default:
			}
		}
	}()

	// Worker function
	workerFunction := func(workerRule *regexp.Regexp, workerReader *io.PipeReader) {
		if workerRule.MatchReader(bufio.NewReader(workerReader)) {
			// Mark us found
			select {
			case foundMatch <- true:
			case <-workerContext.Done():
			}
		} else {
			// Mark us done
			select {
			case finishedWorkers <- true:
			case <-workerContext.Done():
			}
		}
	}

	// Spawn worker for each rule
	for i, rule := range *rules {
		go workerFunction(rule, workerReaders[i])
	}

	finishedWorkersCount := 0

	// Wait to see if we found a match
	for {
		select {
		case <-ctx.Done():
			cancelWorkers()
			return false
		case <-foundMatch:
			cancelWorkers()
			return true
		case <-finishedWorkers: // a go routines finished
			finishedWorkersCount++
			if finishedWorkersCount == rulesLen {
				// All go routines finished
				cancelWorkers()
				return false
			}
		}
	}
}
