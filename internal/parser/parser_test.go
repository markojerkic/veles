package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatterns(t *testing.T) {

	testCases := []struct {
		path    string
		pattern string
		want    bool
	}{
		{"cmd/main.go", "./cmd/main.go", true},
		{"cmd/main.go", "./cmd/*.go", true},
		{"cmd/main.go", "cmd/*.go", true},
		{"cmd/main.go", "./*/*.go", true},
		{"cmd/main.go", "./**/*.go", true},
		{"cmd/main.go", "**/*.go", true},
		{"cmd/main.go", "./**/test.go", false},
		{"cmd/main.go", "./**/main.*", true},
	}

	for _, testCase := range testCases {
		pattern, err := NewPattern(testCase.pattern)
		assert.NoError(t, err)
		assert.Equal(t, testCase.want, pattern.Matches(testCase.path))
	}

}
