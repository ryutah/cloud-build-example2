package server_test

import (
	"testing"

	"gopkg.in/yaml.v2"
)

func TestSample(t *testing.T) {
	_, err := yaml.Marshal(map[string]string{
		"foo": "bar",
	})
	if err != nil {
		t.Fatal(err)
	}
}
