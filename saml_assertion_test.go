package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"encoding/base64"
)

func loadFixture(name string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join("fixtures", name))
}

func s(s string) *string {
	return &s
}

func TestParseAssertion(t *testing.T) {
	b, _ := loadFixture("samlassertion.xml")
	assertion := NewSamlAssertion(base64.StdEncoding.EncodeToString(b))

	assert.NotNil(t, assertion.Parse())

	assert.NotNil(t, assertion.GetProvider("arn:aws:iam::11111111111111:role/some-role"))
	assert.Equal(t, "arn:aws:iam::11111111111111:saml-provider/some-provider", *assertion.GetProvider("arn:aws:iam::11111111111111:role/some-role"))
	assert.NotNil(t, assertion.GetProvider("arn:aws:iam::22222222222222:role/some-role"))
	assert.Equal(t, "arn:aws:iam::22222222222222:saml-provider/some-provider", *assertion.GetProvider("arn:aws:iam::22222222222222:role/some-role"))
	assert.NotNil(t, assertion.GetProvider("arn:aws:iam::22222222222222:role/some-other-role"))
	assert.Equal(t, "arn:aws:iam::22222222222222:saml-provider/some-provider", *assertion.GetProvider("arn:aws:iam::22222222222222:role/some-other-role"))
}

