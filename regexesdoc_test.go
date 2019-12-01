package multiregex

import (
	"context"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

func ExampleRuleSet_MatchesRules() {
	checkMe := []byte("This is a string of text")

	// Check to make sure it matches
	rules := RuleSet{
		regexp.MustCompile(`random text to test for`),
		regexp.MustCompile(`random text to test for two`),
		regexp.MustCompile(`string o`),
		regexp.MustCompile(`random text to test for three`),
	}
	fmt.Println(rules.MatchesRules(checkMe))

	// Output: true
}

func ExampleRuleSet_GetMatchedRules() {
	checkMe := []byte("This is a string of text")

	// Check to make sure it matches
	rules := RuleSet{
		regexp.MustCompile(`random text to test for`),
		regexp.MustCompile(`random text to test for two`),
		regexp.MustCompile(`string o`),
		regexp.MustCompile(`random text to test for three`),
	}
	fmt.Println(rules.GetMatchedRules(checkMe))

	// Output: [string o]
}

func ExampleRuleSet_MatchesRulesReader() {
	checkMe := ioutil.NopCloser(strings.NewReader("This is a string of text"))

	// Check to make sure it matches
	rules := RuleSet{
		regexp.MustCompile(`random text to test for`),
		regexp.MustCompile(`random text to test for two`),
		regexp.MustCompile(`string o`),
		regexp.MustCompile(`random text to test for three`),
	}
	fmt.Println(rules.MatchesRulesReader(context.Background(), checkMe))

	// Output: true
}

func ExampleRuleSet_GetMatchedRulesReader() {
	checkMe := ioutil.NopCloser(strings.NewReader("This is a string of text"))

	// Check to make sure it matches
	rules := RuleSet{
		regexp.MustCompile(`random text to test for`),
		regexp.MustCompile(`random text to test for two`),
		regexp.MustCompile(`string o`),
		regexp.MustCompile(`random text to test for three`),
	}
	fmt.Println(rules.GetMatchedRulesReader(context.Background(), checkMe))

	// Output: [string o]
}
