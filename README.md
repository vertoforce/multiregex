# Regex Machine

This library provides some useful regex helpers when you have a **set of regexes** you need to scan across data.

In this library we call a set of regex rules (`[]*regexp.Regexp`) a `RuleSet`.

## Usage

Currently there are four functions available.

`func (rules RuleSet) GetMatchedRules(data []byte) RuleSet {}`
`func (rules RuleSet) MatchesRules(data []byte) bool {}`

- These function returns matched rules or returns true (ASAP) on the provided data

`func (rules RuleSet) GetMatchedRulesReader(ctx context.Context, reader io.ReadCloser) RuleSet {}`
`func (rules RuleSet) MatchesRulesReader(ctx context.Context, reader io.ReadCloser) bool {}`

- These function returns matched rules or returns true (ASAP) on the reader
- This function opens multiple go routines to scan the input stream concurrently against all rules
- This function will _read_ all data in the reader.
