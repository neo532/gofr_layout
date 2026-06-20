// generate by wireGenerate.sh with '^func New' in on package
package base

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewDatabaseDefault,
	NewRedisLock,
)
