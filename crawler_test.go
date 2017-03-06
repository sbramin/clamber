package main

import "testing"

func TestPageType(t *testing.T) {
	var p pageType
	p.URL = "http://sbramin.com"
	if p.URL != "http://sbramin.com" {
		t.Error("Couldn't store a url in pageType p")
	}
}
