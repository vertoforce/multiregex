package regexmachine

import (
	"context"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"
)

func TestGetMatchedRules(t *testing.T) {

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
