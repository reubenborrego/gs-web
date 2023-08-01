package web

import (
	"regexp"
)

type Matcher interface {
	Match(string) bool
	String() string
	Process(string, string) string
}

type MatchString struct {
	match string
}

func NewMatchString(match string) *MatchString {
	return &MatchString{match: match}
}

func (matchString *MatchString) Match(match string) bool {
	return matchString.match == match
}

func (matchString *MatchString) String() string {
	return matchString.match
}

func (matchString *MatchString) Process(segment, path string) string {
	return path
}

type MatchRegex struct {
	regex *regexp.Regexp
}

func NewMatchRegex(match string) *MatchRegex {
	regex := regexp.MustCompile(match)
	return &MatchRegex{regex: regex}
}

func (matchRegex *MatchRegex) Match(match string) bool {
	return matchRegex.regex.MatchString(match)
}

func (matchRegex *MatchRegex) String() string {
	return matchRegex.regex.String()
}

func (matchRegex *MatchRegex) Process(segment, path string) string {
	return matchRegex.regex.ReplaceAllString(segment, path)
}
