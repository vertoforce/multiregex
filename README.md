# Multi Regex

This library provides some useful regex helpers when you have a **set of regexes** you need to scan across data.

In this library we call a set of regex rules (`[]*regexp.Regexp`) a `RuleSet`.

## Usage

Currently there are four functions available.

```go
func (rules RuleSet) GetMatchedRules(data []byte) RuleSet {}
func (rules RuleSet) MatchesRules(data []byte) bool {}
```

- These function returns matched rules or returns true (ASAP) on the provided data

```go
func (rules RuleSet) GetMatchedRulesReader(ctx context.Context, reader io.ReadCloser) RuleSet {}
func (rules RuleSet) MatchesRulesReader(ctx context.Context, reader io.ReadCloser) bool {}
```

- These function returns matched rules or returns true (ASAP) on the reader
- These functions open multiple go routines to scan the input stream concurrently against all rules
- These functions will _read all data_ in the reader.

## Example

```go
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
```
