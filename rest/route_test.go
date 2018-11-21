package rest_test

import (
	"testing"

	"github.com/spec-tacles/spectacles/rest"
)

var routes = map[string]string{
	"/guilds/1/members/1/roles": "/guilds/1/members/:id/roles",
	"/channels/1":               "/channels/1",
	"/users":                    "/users",
	"/users/1":                  "/users/:id",
}

func TestMakeRoute(t *testing.T) {
	for path, route := range routes {
		if madeRoute := rest.MakeRoute(path); madeRoute != route {
			t.Errorf("Constructed route '%s' not equal to expected route '%s'", madeRoute, route)
			return
		}
	}
}
