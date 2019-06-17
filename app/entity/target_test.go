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

func TestTarget_MakeURL(t *testing.T) {
	// GIVEN a target with an URL template
	target, err := NewTarget("test1", "", "https://ahost.com:8080{{.Path}}")
	assert.NoError(t, err)

	// WHEN an URL is requested with a Message with path
	msg := Message{Path: "/aNice/Path/"}
	targetURL, err := target.MakeURL(&msg)

	// THEN the URL is generated
	assert.NoError(t, err)
	assert.Equal(t, "https://ahost.com:8080/aNice/Path/", targetURL.String())
}

func TestTarget_GetIDFor_NoMatch(t *testing.T) {
	// GIVEN A target with URL template
	target, _ := NewTarget("test1", "/thisTarget", "")

	// WHEN A path is not matched
	msg := &Message{Path:"/noMatch/expected/here"}
	id := target.GetIDFor(msg)

	// THEN the ID is all the path
	assert.Equal(t, msg.Path, id)
}

func TestTarget_GetIDFor_MatchWithoutCapturingGroup(t *testing.T) {
	// GIVEN A target with URL template
	target, _ := NewTarget("test1", "/thisTarget", "")

	// WHEN A path is not matched
	msg := &Message{Path:"/thisTarget/expected/here"}
	id := target.GetIDFor(msg)

	// THEN the ID is all the path
	assert.Equal(t, "/thisTarget", id)
}

func TestTarget_GetIDFor_MatchWithCapturingGroup(t *testing.T) {
	// GIVEN A target with URL template
	target, err := NewTarget("test1", "/thisTarget/([0-9]+)/p", "")
	assert.NoError(t, err)

	// WHEN A path is not matched
	msg := &Message{Path:"/thisTarget/345873489534/p"}
	id := target.GetIDFor(msg)

	// THEN the ID is all the path
	assert.Equal(t, "345873489534", id)
}

