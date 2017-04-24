package main

import (
	"encoding/json"
	"testing"
)

type testCase struct {
	input    []byte
	expected bool
	url      string
	children string
	atype    string
	aurl     string
}

var testCases = []testCase{
	{[]byte(`{"url":"https://sbramin.com/", "children":["https://sbramin.com/about/"],"assets":[{"url":"https://sbramin.com/css/base.css","type":"css"}]}`), true, "https://sbramin.com/", "https://sbramin.com/about/", "css", "https://sbramin.com/css/base.css"},
	{[]byte(`{"crap json"}}`), false, "", "", "", ""},
}

func TestPageReview(t *testing.T) {
	var pt page
	actual := true
	for _, test := range testCases {
		err := json.Unmarshal(test.input, &pt)
		if err != nil {
			actual = false
		}
		if actual != test.expected {
			t.Errorf("Did not match, expected %s got %s", actual, test.expected)
		}

		if actual == false {
			continue
		}

		if pt.URL != test.url {
			t.Errorf("URL does not match, got %s, expected %s", pt.URL, test.url)
		}
		if pt.Children[0] != test.children {
			t.Errorf("Children do not match, got %s, expected %s", pt.Children[0], test.children)
		}
		if pt.Assets[0].Type != test.atype {
			t.Errorf("URL does not match, got %s, expected %s", pt.Assets[0].Type, test.atype)
		}
		if pt.Assets[0].Type != test.atype {
			t.Errorf("Asset URL does not match, got %s, expected %s", pt.Assets[0].Type, test.atype)
		}

	}
}
