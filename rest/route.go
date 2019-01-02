package rest

import (
	"path"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	// PathSep is a string representing the delimiting character in a URL path
	PathSep = "/"

	// IDNotation is a string representing an ID that has been replaced in the formation of a ratelimit route
	IDNotation = PathSep + ":id"
)

// MakeRoute makes a ratelimit route given a path
func MakeRoute(p string) (route string) {
	var params = strings.Split(p[1:], PathSep)

	if len(params) == 0 {
		return ""
	}

	if params[0] == "channels" || params[0] == "guilds" || params[0] == "webhooks" {
		// channels, guilds, and webooks are considered "primary IDs"
		if len(params) == 1 {
			return params[0]
		}

		var route = path.Join(params[0], params[1])
		return buildRoute(params, route, 2)
	}

	return buildRoute(params, "", 0)
}

func buildRoute(params []string, route string, i int) string {
	for ; i < len(params); i++ {
		// if the first character of the parameter is a number, consider it to be an ID
		r, _ := utf8.DecodeRuneInString(params[i])
		if unicode.IsDigit(r) {
			route += IDNotation
		} else {
			route += path.Join(params[i])
		}
	}

	return route
}
