# Regex Machine

This library provides some useful regex helpers when you have a **set of regexes** you need to scan across data.

In this library we call a set of regex rules (`[]*regexp.Regexp`) a `RuleSet`.

## Usage

Currently there are two functions available.

`func (rules RuleSet) GetMatchedRules(data []byte) RuleSet {}`

- This function returns matched rules on the provided data

`func (rules RuleSet) GetMatchedRulesReader(ctx context.Context, reader io.ReadCloser) RuleSet {}`

- This function returns matches rules in the reader.
- This function opens multiple go routines to scan the input stream concurrently against all rules
- This function will _read_ all data in the reader.
