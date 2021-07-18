package appcast

import (
	"regexp"
)

type RewriteFunc func(string) string

func RegexRewrite(regex, replace string) RewriteFunc {
	r := regexp.MustCompile(regex)
	return func(s string) string { return r.ReplaceAllString(s, replace) }
}
