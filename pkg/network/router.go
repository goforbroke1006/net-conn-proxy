package network

import (
	"net"
	"regexp"
)

type HandleFunc func(payload []byte, w net.Conn)

type Router struct {
	rules []routerRule
}

func (r *Router) HandleFunc(pattern string, fn HandleFunc) {
	r.rules = append(r.rules, routerRule{
		patternRe: regexp.MustCompile(pattern),
		fn:        fn,
	})
}

func (r *Router) getMatch(payload []byte) HandleFunc {
	for idx := 0; idx < len(r.rules); idx++ {
		if r.rules[idx].patternRe.Match(payload) {
			return r.rules[idx].fn
		}
	}

	return func(_ []byte, _ net.Conn) {}
}

type routerRule struct {
	patternRe *regexp.Regexp
	fn        HandleFunc
}
