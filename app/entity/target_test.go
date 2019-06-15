package entity

import (
	"testing"
	"github.com/stretchr/testify/assert"
)
func TestTarget_Match_Match(t *testing.T) {
	// GIVEN a target with URL perfix expression
	target, err := NewTarget("test1", "/thisTarget", "")
	assert.NoError(t, err)

	// WHEN try to match a handled path
	match := target.Match("/thisTarget/withOther/paths")

	// THEN it matches
	assert.True(t, match)
}

func TestTarget_Match_NoMatch(t *testing.T) {
	// GIVEN a target with URL perfix expression
	target, err := NewTarget("test1", "/thisTarget", "")
	assert.NoError(t, err)

	// WHEN try to match a not handled path
	match := target.Match("/thisIsOtherTarget/withOther/paths")

	// THEN it does not matches
	assert.False(t, match)
}
