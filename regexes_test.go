package regexmachine

import (
	"bytes"
	"context"
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
	matchedRules := rules.GetMatchedRulesReader(context.Background(), strings.NewReader(checkMe))
	if len(matchedRules) == 0 {
		t.Errorf("Should have matched")
	}

	var buf bytes.Buffer

	// Write 1MB of "A" to the buffer
	// Then "text to match"
	for i := 0; i < 1024*1024; i++ {
		buf.Write([]byte("A"))
	}
	buf.Write([]byte("matchme!"))

	rules = RuleSet{
		regexp.MustCompile(`AB`),
		regexp.MustCompile(`matchme`),
	}

	// Should match
	matchedRules = rules.GetMatchedRulesReader(context.Background(), &buf)
	if len(matchedRules) == 0 {
		t.Errorf("Should have matched")
	}
}
