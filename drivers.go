package gloader

// Note: All built-in drivers are registered in the init() function of the driver package.
// and for registration mechanism to work, the driver packages must be imported here.
// nolint:revive // ignore unused import warning.
import (
	_ "github.com/mohammadv184/gloader/driver/cockroach"
	_ "github.com/mohammadv184/gloader/driver/mysql"
)
