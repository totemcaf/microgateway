package entity

import (
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
func NewTarget(id string, receiveURLPatternStr string, targetURLStr string) (*Target, error) {
	receiveURLPatern, err := regexp.Compile(targetURLStr)
	if err != nil {
		return nil, err
	}

	urlTemplate, err := template.New("targetUrl").Parse(receiveURLPatternStr)
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
	return t.receiveURLPatern.MatchString(path)
}
