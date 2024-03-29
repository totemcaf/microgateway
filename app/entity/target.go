package entity

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"text/template"
)

// Target contains configuration to generate the destination for the message
type Target struct {
	// ID  identifies this target
	id               string
	urlTemplate      *template.Template
	receiveURLPatern *regexp.Regexp
}

// NewTarget builds a new Target.
// The receiveURLPatternStr should be a valid go Template.
// In the template the provided variables are:
//		.method	the called method
// 		.path	the called path (includes the query string)
func NewTarget(id string, receiveURLPatternStr string, targetURLTemplateStr string) (*Target, error) {
	receiveURLPatern, err := regexp.Compile(receiveURLPatternStr)
	if err != nil {
		return nil, err
	}

	urlTemplate, err := template.New("targetUrl").Parse(targetURLTemplateStr)
	if err != nil {

		return nil, err
	}

	return &Target{id, urlTemplate, receiveURLPatern}, nil
}

// MakeURL returnd the URL to access that target for the given message
func (t *Target) MakeURL(m *Message) (*url.URL, error) {

	var buff strings.Builder

	if err := t.urlTemplate.Execute(&buff, m); err != nil {
		return nil, err
	}

	return url.Parse(buff.String())
}

// Match checks this target agains the given URL path. Return true if the path is for this target
func (t *Target) Match(path string) bool {
	fmt.Print("->")
	fmt.Println(t.receiveURLPatern.String())
	return t.receiveURLPatern.MatchString(path)
}

// GetIDFor extracts the message ID from the message.
// This versions extracts the ID from the message "path" using the regular expression.
func (t *Target) GetIDFor(m *Message) string {
	groups := t.receiveURLPatern.FindStringSubmatch(m.Path)

	switch len(groups) {
	case 0:
		// No match, defaults to all the path
		return m.Path
	case 1:
		// All the path is matched
		return groups[0]
	default:
		// At least 1 subgroup is found, use it
		return groups[1]
	}
}
