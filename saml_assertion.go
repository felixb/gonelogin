package main

import (
	"encoding/base64"
	"regexp"
	"strings"
)

type SamlAssertion struct {
	Assertion *string
	parsed    map[string]*string
}

func NewSamlAssertion(assertion string) *SamlAssertion {
	return &SamlAssertion{
		Assertion: &assertion,
	}
}

func (a *SamlAssertion) Parse() map[string]*string {
	if a.parsed != nil {
		return a.parsed
	}

	b, err := base64.StdEncoding.DecodeString(*a.Assertion)
	if err != nil {
		// TODO something sane
		return nil
	}

	re := regexp.MustCompile(`arn:aws:iam::[0-9]+:role/[^,]+,arn:aws:iam::[0-9]+:saml-provider/[^\s<]+`)
	f := re.FindAll(b, -1)
	if f == nil {
		// TODO something sane?
		return nil
	}

	a.parsed = make(map[string]*string)
	for _, match := range f {
		s := string(match)
		params := strings.Split(s, ",")
		a.parsed[params[0]] = &params[1]
	}

	return a.parsed
}

func (a *SamlAssertion) GetProvider(roleArn string) *string {
	a.Parse()
	if a.parsed == nil {
		return nil
	}
	return a.parsed[roleArn]
}
