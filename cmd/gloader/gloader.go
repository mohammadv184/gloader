package gloader

import (
	"database/sql"
	"sync"
)

var it *string

func main() {
	sql.Drivers()

	sync.NewCond(&sync.Mutex{})
}
