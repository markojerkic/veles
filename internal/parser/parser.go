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

func isContainingWildcard(part string) ([]string, bool) {

	wildcardParts := strings.Split(part, "*")
	numOfWildcards := len(wildcardParts)

	if numOfWildcards > 1 && part != string(WILDCARD_FORWARD_DIR) {
		panic("Path pattern my contain ** only in context of ./**/*.go")
	}

	return wildcardParts, numOfWildcards == 1

}

func (self *Pattern) Matches(path string) bool {

	pathParts := strings.Split(path, "/")

	pathIndex := len(pathParts)
	patternIndex := len(self.parts)

	wildcardNextTarget := ""

	for pathIndex > 0 && patternIndex > 0 {

		if parts, ok := isContainingWildcard(pathParts[pathIndex]); ok {
			if !strings.HasPrefix(pathParts[pathIndex], parts[0]) {
				return false
			}
		}

		if pathParts[pathIndex] == string(WILDCARD_DIR) {
			pathIndex--
			patternIndex--
			continue
		}

		if pathParts[pathIndex] == string(WILDCARD_FORWARD_DIR) {
			pathIndex--
			if pathIndex == 0 {
				return false
			}
			wildcardNextTarget = self.parts[pathIndex-1]
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
