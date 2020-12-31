package gorr

import (
	"fmt"
	"net/http"
)

func Static(t string) NodeHeader {
	return NodeHeader{
		title: t,
		match: Matches(t),
	}
}
func Parameter(t string, m Matcher) NodeHeader {
	return NodeHeader{
		title: fmt.Sprintf("<%s>", t),
		match: m,
	}
}
func Matches(expected string) Matcher {
	return func(v string) bool {
		return expected == v
	}
}
func MatchesOneOf(samples []string) Matcher {
	return func(v string) bool {
		for _, s := range samples {
			if s == v {
				return true
			}
		}
		return false
	}
}
func StringToMethod(s string) (m HTTPMethod) {
	switch s {
	case "GET":
		m = GET
	case "HEAD":
		m = HEAD
	case "POST":
		m = POST
	case "PUT":
		m = PUT
	case "DELETE":
		m = DELETE
	case "CONNECT":
		m = CONNECT
	case "OPTIONS":
		m = OPTIONS
	case "TRACE":
		m = TRACE
	case "PATCH":
		m = PATCH
	default:
		m = WRONG_METHOD
	}
	return
}
func DumbHook(w http.ResponseWriter, r *http.Request) {
	// do nothing
}
