package server

import (
	"maps"

	gofrscript "github.com/neo532/gofr/transport/script"
	"github.com/neo532/gofr_layout/internal/service/script"
)

// NewScriptServer creates a ScriptServer from a route map.
func NewScriptServer(
	user *script.UserScript,
	user1 *script.User1Script,
) *gofrscript.Server {

	routes := make(map[string]gofrscript.Func, 10)

	maps.Copy(routes, gofrscript.Discover(user))
	maps.Copy(routes, gofrscript.Discover(user1))

	return gofrscript.New(routes)
}
