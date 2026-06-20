//go:build generate

package generate

import (
	_ "github.com/google/wire/cmd/wire"
)

//go:generate make config
