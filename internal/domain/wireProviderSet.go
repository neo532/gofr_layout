// generate by wireGenerate.sh with '^func New' in on package
package domain

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUserDomain,
)
