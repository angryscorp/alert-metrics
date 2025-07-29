package buildinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildInfo_String_WithValues(t *testing.T) {
	buildInfo := New("1.0.0", "2023-01-01", "abc123")
	result := buildInfo.String()

	expected := "Build version: 1.0.0\nBuild date: 2023-01-01\nBuild commit: abc123\n"
	assert.Equal(t, expected, result)
}

func TestBuildInfo_String_WithEmptyValues(t *testing.T) {
	buildInfo := New("", "", "")
	result := buildInfo.String()

	expected := "Build version: N/A\nBuild date: N/A\nBuild commit: N/A\n"
	assert.Equal(t, expected, result)
}

func TestBuildInfo_String_WithPartialValues(t *testing.T) {
	buildInfo := New("1.0.0", "", "abc123")
	result := buildInfo.String()

	expected := "Build version: 1.0.0\nBuild date: N/A\nBuild commit: abc123\n"
	assert.Equal(t, expected, result)
}
