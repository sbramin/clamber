package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"urlIndex",
		"GET",
		"/urls",
		clamberIndex,
	},
	Route{
		"urlShow",
		"GET",
		"/urls/{url}",
		clamberShow,
	},
	Route{
		"urlCreate",
		"POST",
		"/urls",
		clamber,
	},
}
