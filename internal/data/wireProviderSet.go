// generate by wireGenerate.sh with '^func New' in on package
package data

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewTransactionDefaultRepo,
	NewUserRepo,
)
