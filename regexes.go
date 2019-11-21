// Package regexmachine helps with processing regex rules of streams of data and looking for rule matches
package regexmachine

import (
	"bufio"
	"context"
	"io"
	"regexp"
)

// RuleSet A set of regex rules
type RuleSet []*regexp.Regexp

// Public regex sets
var (
	Email = regexp.MustCompile(`[A-Za-z0-9_.]+((\ ?(\[|\()?\ ?@\ ?(\)|\])?\ ?)|(\ ?(\[|\()\ ?[aA][tT]\ ?(\)|\])\ ?))[0-9a-z.-]+`)

	// DefaultSet
	DefaultRules = RuleSet{Email}
	MatchAll     = RuleSet{regexp.MustCompile(`.*`)}
)

// -- Functions on RuleSet --

// GetMatchedRules Given bytes return all regexes that match
func (rules RuleSet) GetMatchedRules(data []byte) RuleSet {
	matched := []*regexp.Regexp{}
	for _, rule := range rules {
		if rule.Match(data) {
			matched = append(matched, rule)
		}
	}

	return matched
}

// GetMatchedRulesReader Given a reader, return any rule matches in the stream.  Will read ENTIRE READER
// Spawns multiple go routines to check each rule
// Use limit reader to prevent reading forever
func (rules RuleSet) GetMatchedRulesReader(ctx context.Context, reader io.Reader) RuleSet {
	matchedRules := make(chan *regexp.Regexp)
	finishedWorkers := make(chan bool)

	// We need to duplicate the reader stream for each worker
	// Create reader and writer for each worker thread
	workerWriters := []*io.PipeWriter{}
	workerReaders := []*io.PipeReader{}
	for range rules {
		r, w := io.Pipe()
		workerWriters = append(workerWriters, w)
		workerReaders = append(workerReaders, r)
	}

	ctxWorkers, cancelWorkers := context.WithCancel(ctx)

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
			case <-ctxWorkers.Done():
				return
			default:
			}
		}
	}()

	// Worker function
	workerFunction := func(workerRule *regexp.Regexp, workerReader *io.PipeReader) {
		if workerRule.MatchReader(bufio.NewReader(workerReader)) {
			// Mark this rule as matched
			select {
			case matchedRules <- workerRule:
			case <-ctxWorkers.Done():
				return
			}
		}
		// Mark us done
		select {
		case finishedWorkers <- true:
		case <-ctxWorkers.Done():
			return
		}
	}

	// Spawn worker for each rule
	for i, rule := range rules {
		go workerFunction(rule, workerReaders[i])
	}

	// Routine to capture all matches
	matches := RuleSet{}
	go func() {
		for match := range matchedRules {
			matches = append(matches, match)
		}
	}()

	// Wait for threads to finish
	finishedWorkersCount := 0
	for finishedWorkersCount != len(rules) {
		select {
		case <-ctx.Done():
			cancelWorkers()
			return nil
		case <-finishedWorkers: // a go routine finished
			finishedWorkersCount++
		}
	}

	close(matchedRules)
	cancelWorkers()

	// Return found matches
	return matches
}
