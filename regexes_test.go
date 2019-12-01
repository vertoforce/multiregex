package multiregex

import (
	"context"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"
)

func TestMatchesRules(t *testing.T) {
	checkMe := []byte("This is a string of text")

	// Check to make sure it matches
	rules := RuleSet{
		regexp.MustCompile(`random text to test for`),
		regexp.MustCompile(`random text to test for two`),
		regexp.MustCompile(`string o`),
		regexp.MustCompile(`random text to test for three`),
	}
	if !rules.MatchesRules(checkMe) {
		t.Errorf("Should have matched")
	}

	// Check to make sure it does not match
	rules2 := RuleSet{
		regexp.MustCompile(`random text to test for`),
		regexp.MustCompile(`random text to test for two`),
		regexp.MustCompile(`random text to test for three`),
	}
	if rules2.MatchesRules(checkMe) {
		t.Errorf("Should not have matched")
	}
}

func TestMatchesRulesReader(t *testing.T) {
	checkMe := ioutil.NopCloser(strings.NewReader("This is a string of text"))

	// Check to make sure it matches
	rules := RuleSet{
		regexp.MustCompile(`random text to test for`),
		regexp.MustCompile(`random text to test for two`),
		regexp.MustCompile(`string o`),
		regexp.MustCompile(`random text to test for three`),
	}
	if !rules.MatchesRulesReader(context.Background(), checkMe) {
		t.Errorf("Should have matched")
	}

	// Check to make sure it does not match
	rules2 := RuleSet{
		regexp.MustCompile(`random text to test for`),
		regexp.MustCompile(`random text to test for two`),
		regexp.MustCompile(`random text to test for three`),
	}
	if rules2.MatchesRulesReader(context.Background(), checkMe) {
		t.Errorf("Should not have matched")
	}

	// Check for any lingering go routines
	// time.Sleep(time.Second * 2)
}

func TestGetMatchedRules(t *testing.T) {
	checkMe := "This is a string of text"

	rules := RuleSet{
		regexp.MustCompile(`random text to test for`),
		regexp.MustCompile(`random text to test for two`),
		regexp.MustCompile(`string o`),
		regexp.MustCompile(`random text to test for three`),
	}

	matchedRules := rules.GetMatchedRules([]byte(checkMe))
	if len(matchedRules) != 1 {
		t.Errorf("Did not match correct rules")
	}
}

func TestGetMatchedRulesReader(t *testing.T) {
	checkMe := "This is a string of text"

	rules := RuleSet{
		regexp.MustCompile(`random text to test for`),
		regexp.MustCompile(`random text to test for two`),
		regexp.MustCompile(`string o`),
		regexp.MustCompile(`random text to test for three`),
	}

	// Should match
	matchedRules := rules.GetMatchedRulesReader(context.Background(), ioutil.NopCloser(strings.NewReader(checkMe)))
	if len(matchedRules) == 0 {
		t.Errorf("Should have matched")
	}

	pipeReader, pipeWriter := io.Pipe()

	rules = RuleSet{
		regexp.MustCompile(`AB`),
		regexp.MustCompile(`matchme`),
	}

	// Write 10KB of "A" to the buffer
	// Then "text to match"
	go func() {
		for i := 0; i < 1024*10; i++ {
			pipeWriter.Write([]byte("AAAAAAAAAA"))
		}
		pipeWriter.Write([]byte("matchme!"))
		pipeWriter.Close()
	}()

	// Should match
	matchedRules = rules.GetMatchedRulesReader(context.Background(), pipeReader)
	if len(matchedRules) == 0 {
		t.Errorf("Should have matched")
	}
}
