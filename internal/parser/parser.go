package parser

import (
	"fmt"
	"strings"
)

type Pattern struct {
	pattern string

	parts []string
}

type PatternToken string

const (
	CURRENT_DIR          PatternToken = "."
	PARENT_DIR           PatternToken = ".."
	WILDCARD_DIR         PatternToken = "*"
	WILDCARD_FORWARD_DIR PatternToken = "**"
)

func (self *Pattern) isContainingWildcard(i int) ([]string, bool) {

	wildcardParts := strings.Split(self.parts[i], "*")
	numOfWildcards := strings.Count(self.parts[i], "*")

	if numOfWildcards > 1 && self.parts[i] != string(WILDCARD_FORWARD_DIR) {
		panic("Path pattern my contain ** only in context of ./**/*.go")
	}

	return wildcardParts, numOfWildcards == 1
}

func (self *Pattern) Matches(path string) bool {

	pathParts := strings.Split(path, "/")

	pathIndex := len(pathParts) - 1
	patternIndex := len(self.parts) - 1

	wildcardNextTarget := ""

	for pathIndex >= 0 && patternIndex >= 0 {

		if parts, ok := self.isContainingWildcard(patternIndex); ok {
			hasPrefix := strings.HasPrefix(pathParts[pathIndex], parts[0])
			hasSuffix := strings.HasSuffix(pathParts[pathIndex], parts[1])
			if !hasPrefix || !hasSuffix {
				return false
			}
			pathIndex--
			patternIndex--
			continue
		}

		if self.parts[patternIndex] == string(WILDCARD_DIR) {
			pathIndex--
			patternIndex--
			continue
		}

		if self.parts[patternIndex] == string(WILDCARD_FORWARD_DIR) {
			pathIndex--
			if pathIndex <= 0 {
				return false
			}
			wildcardNextTarget = self.parts[patternIndex-1]
			continue
		}

		if wildcardNextTarget != "" {
			if pathParts[pathIndex] == wildcardNextTarget {
				patternIndex--
				wildcardNextTarget = ""
			}
			pathIndex--
			continue
		}

		if pathParts[pathIndex] == string(CURRENT_DIR) {
			pathIndex--
			continue
		}

		if pathParts[pathIndex] == string(PARENT_DIR) {
			panic("Nisam još ovo skužio, triba razmisliti")
		}

		if pathParts[pathIndex] == self.parts[patternIndex] {
			pathIndex--
			patternIndex--
			continue
		}
		return false

	}

	return true
}

func NewPattern(pattern string) (Pattern, error) {

	parts := strings.Split(pattern, "/")

	if len(parts) == 0 {
		return Pattern{}, fmt.Errorf("Pattern is empty")
	}

	return Pattern{pattern, parts}, nil
}
